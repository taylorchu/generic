// Package generic generates package with type replacements.
package generic

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"github.com/taylorchu/generic/importer"

	"golang.org/x/tools/go/ast/astutil"
)

// rewritePkgName sets current package name.
func rewritePkgName(node *ast.File, pkgName string) {
	node.Name.Name = pkgName
}

// rewriteIdent converts TypeXXX to its replacement defined in typeMap.
func rewriteIdent(node *ast.File, typeMap map[string]Target, fset *token.FileSet) {
	var used []string
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Obj == nil || x.Obj.Kind != ast.Typ {
				return false
			}
			to, ok := typeMap[x.Name]
			if !ok {
				return false
			}
			x.Name = to.Ident

			if to.Import == "" {
				return false
			}
			var found bool
			for _, im := range used {
				if im == to.Import {
					found = true
					break
				}
			}
			if !found {
				used = append(used, to.Import)
			}
			return false
		}
		return true
	})
	for _, im := range used {
		astutil.AddImport(fset, node, im)
	}
}

// removeTypeDecl removes type declarations defined in typeMap.
func removeTypeDecl(node *ast.File, typeMap map[string]Target) {
	for i := len(node.Decls) - 1; i >= 0; i-- {
		genDecl, ok := node.Decls[i].(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		var remove bool
		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)

			_, ok = typeMap[typeSpec.Name.Name]
			if !ok {
				continue
			}

			_, ok = typeSpec.Type.(*ast.Ident)
			if !ok {
				continue
			}
			remove = true
			break
		}
		if remove {
			node.Decls = append(node.Decls[:i], node.Decls[i+1:]...)
		}
	}
}

// findTypeDecl finds type and related declarations.
func findTypeDecl(node *ast.File) (ret []ast.Decl) {
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec := spec.(*ast.TypeSpec)

			// Replace a complex declaration with a dummy idenifier.
			//
			// It seems simpler to check whether a type is defined.
			typeSpec.Type = &ast.Ident{
				Name: "uint32",
			}
		}

		ret = append(ret, decl)
	}
	return
}

// rewriteTopLevelIdent adds a prefix to top-level identifiers and their uses.
//
// This prevents name conflicts when a package is rewritten to $PWD.
func rewriteTopLevelIdent(nodes map[string]*ast.File, prefix string, typeMap map[string]Target) {
	prefixIdent := func(name string) string {
		if name == "_" {
			// skip unnamed
			return "_"
		}
		return lintName(fmt.Sprintf("%s_%s", prefix, name))
	}

	declMap := make(map[interface{}]string)

	for _, node := range nodes {
		for _, decl := range node.Decls {
			switch decl := decl.(type) {
			case *ast.FuncDecl:
				if decl.Recv != nil {
					continue
				}
				decl.Name.Name = prefixIdent(decl.Name.Name)
				declMap[decl] = decl.Name.Name
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						obj := spec.Name.Obj
						if obj != nil && obj.Kind == ast.Typ {
							if to, ok := typeMap[obj.Name]; ok && spec.Name.Name == to.Ident {
								// If this identifier is already rewritten before, we don't need to prefix it.
								continue
							}
						}
						spec.Name.Name = prefixIdent(spec.Name.Name)
						declMap[spec] = spec.Name.Name
					case *ast.ValueSpec:
						for _, ident := range spec.Names {
							ident.Name = prefixIdent(ident.Name)
							declMap[spec] = ident.Name
						}
					}
				}
			}
		}
	}

	// After top-level identifiers are renamed, find where they are used, and rewrite those.
	for _, node := range nodes {
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				if x.Obj == nil || x.Obj.Decl == nil {
					return false
				}
				name, ok := declMap[x.Obj.Decl]
				if !ok {
					return false
				}
				x.Name = name
				return false
			}
			return true
		})
	}
}

type packageTarget struct {
	SameDir bool
	NewName string
	NewPath string
}

// parsePackageTarget finds where a package can be rewritten.
func parsePackageTarget(path string) (*packageTarget, error) {
	t := new(packageTarget)
	if strings.HasPrefix(path, ".") {
		t.SameDir = true
		t.NewPath = strings.TrimPrefix(path, ".")
		t.NewName = os.Getenv("GOPACKAGE")
		if t.NewName == "" {
			return nil, errors.New("GOPACKAGE cannot be empty")
		}
	} else {
		t.NewPath = path
		t.NewName = filepath.Base(path)
	}

	return t, nil
}

// RewritePackage applies type replacements on a package in $GOPATH, and saves results as a new package in $PWD.
//
// If there is a dir with the same name as newPkgPath, it will first be removed. It is possible to re-run this
// to update a generic package.
func RewritePackage(pkgPath string, newPkgPath string, typeMap map[string]Target) error {
	var err error

	pt, err := parsePackageTarget(newPkgPath)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	files := make(map[string]*ast.File)

	// NOTE: this package that we try to rewrite from should not contain vendor/.
	buildP, err := build.Import(pkgPath, "", 0)
	if err != nil {
		return err
	}
	for _, file := range buildP.GoFiles {
		path := filepath.Join(buildP.Dir, file)
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		files[path] = f
	}

	// Gather ast.File to create ast.Package.
	// ast.NewPackage will try to resolve unresolved identifiers.
	ast.NewPackage(fset, files, nil, nil)

	// Apply AST changes and refresh.
	for _, f := range files {
		rewritePkgName(f, pt.NewName)
		removeTypeDecl(f, typeMap)
		rewriteIdent(f, typeMap, fset)
	}

	if pt.SameDir {
		rewriteTopLevelIdent(files, pt.NewPath, typeMap)
	}

	// AST in dirty state; refresh
	buf := new(bytes.Buffer)
	var tc []*ast.File
	for path, f := range files {
		buf.Reset()
		err = printer.Fprint(buf, fset, f)
		if err != nil {
			return err
		}
		f, err = parser.ParseFile(fset, "", buf, 0)
		if err != nil {
			printer.Fprint(os.Stderr, fset, f)
			return err
		}
		files[path] = f
		tc = append(tc, f)
	}

	outPath := func(path string) string {
		if pt.SameDir {
			return fmt.Sprintf("%s_%s", pt.NewPath, filepath.Base(path))
		}
		return filepath.Join(pt.NewPath, filepath.Base(path))
	}

	// Type-check.
	if pt.SameDir {
		// Also include same-dir files.
		// However, it is silly to add the entire file,
		// because that file might have identifiers from another generic package.
		buildP, err := build.Import(".", ".", 0)
		if err != nil {
			if _, ok := err.(*build.NoGoError); !ok {
				return err
			}
		}
		generated := func(path string) bool {
			for p := range files {
				if outPath(p) == path {
					return true
				}
			}
			return false
		}
		for _, file := range buildP.GoFiles {
			path := filepath.Join(buildP.Dir, file)
			if generated(path) {
				// Allow updating existing generated files.
				continue
			}
			f, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}
			decl := findTypeDecl(f)
			if len(decl) > 0 {
				tc = append(tc, &ast.File{
					Decls: decl,
					Name:  f.Name,
				})
			}
		}
	}
	conf := types.Config{Importer: importer.New()}
	_, err = conf.Check("", fset, tc, nil)
	if err != nil {
		for _, f := range tc {
			printer.Fprint(os.Stderr, fset, f)
		}
		return err
	}

	writeOutput := func() error {
		for path, f := range files {
			// Print ast to file.
			var dest *os.File
			dest, err = os.Create(outPath(path))
			if err != nil {
				return err
			}
			defer dest.Close()

			err = format.Node(dest, fset, f)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if pt.SameDir {
		return writeOutput()
	}

	err = os.RemoveAll(pt.NewPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(pt.NewPath, 0777)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(pt.NewPath)
		}
	}()

	return writeOutput()
}
