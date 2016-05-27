// Package generic generates package with type replacements.
package generic

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if _, ok := typeMap[typeSpec.Name.Name]; ok {
				remove = true
				break
			}
		}
		if remove {
			node.Decls = append(node.Decls[:i], node.Decls[i+1:]...)
		}
	}
}

// rewriteTopLevelIdent adds a prefix to top-level identifiers and their uses.
// For example, XXX will be converted to _prefix_XXX.
//
// This prevents name conflicts when a package is rewritten to PWD.
func rewriteTopLevelIdent(node *ast.File, prefix string) {
	declMap := make(map[interface{}]string)
	for i := len(node.Decls) - 1; i >= 0; i-- {
		switch decl := node.Decls[i].(type) {
		case *ast.FuncDecl:
			if decl.Recv == nil {
				decl.Name.Name = fmt.Sprintf("%s_%s", prefix, decl.Name.Name)
				declMap[decl] = decl.Name.Name
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch x := spec.(type) {
				case *ast.TypeSpec:
					x.Name.Name = fmt.Sprintf("%s_%s", prefix, x.Name.Name)
					declMap[x] = x.Name.Name
				case *ast.ValueSpec:
					for _, ident := range x.Names {
						ident.Name = fmt.Sprintf("%s_%s", prefix, ident.Name)
						declMap[x] = ident.Name
					}
				}
			}
		}
	}

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

// walkSource visits all .go files in a package path except tests.
func walkSource(pkgPath string, sourceFunc func(string) error) error {
	pkgPath = fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkgPath)
	fi, err := ioutil.ReadDir(pkgPath)
	if err != nil {
		return err
	}
	for _, info := range fi {
		if info.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s/%s", pkgPath, info.Name())
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		err = sourceFunc(path)
		if err != nil {
			return err
		}
	}
	return nil
}

type packageTarget struct {
	SameDir bool
	NewName string
	NewPath string
}

func parsePackageTarget(path string) (*packageTarget, error) {
	t := new(packageTarget)
	if strings.HasPrefix(path, ".") {
		t.SameDir = true
		t.NewPath = strings.Replace(strings.TrimPrefix(path, "."), "/", "_", -1)
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

// RewritePackage applies type replacements on a package in GOPATH, and saves results as a new package in $PWD.
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
	err = walkSource(pkgPath, func(path string) error {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		files[path] = f
		return nil
	})
	if err != nil {
		return err
	}

	// Gather ast.File to create ast.Package.
	// ast.NewPackage will try to resolve unresolved identifiers.
	ast.NewPackage(fset, files, nil, nil)

	// Apply AST changes and refresh.
	buf := new(bytes.Buffer)
	var tc []*ast.File
	for path, f := range files {
		rewritePkgName(f, pt.NewName)
		removeTypeDecl(f, typeMap)
		rewriteIdent(f, typeMap, fset)
		if pt.SameDir {
			rewriteTopLevelIdent(f, pt.NewPath)
		}

		// AST in dirty state; refresh
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

	// Type check.
	conf := types.Config{Importer: importer.Default()}
	_, err = conf.Check("", fset, tc, nil)
	if err != nil {
		for _, f := range tc {
			printer.Fprint(os.Stderr, fset, f)
		}
		return err
	}

	if pt.SameDir {
		for path, f := range files {
			// Print ast to file.
			var dest *os.File
			dest, err = os.Create(fmt.Sprintf("%s_%s", pt.NewPath, filepath.Base(path)))
			if err != nil {
				return err
			}
			defer dest.Close()

			err = printer.Fprint(dest, fset, f)
			if err != nil {
				return err
			}
		}
		return nil
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

	for path, f := range files {
		// Print ast to file.
		var dest *os.File
		dest, err = os.Create(fmt.Sprintf("%s/%s", pt.NewPath, filepath.Base(path)))
		if err != nil {
			return err
		}
		defer dest.Close()

		err = printer.Fprint(dest, fset, f)
		if err != nil {
			return err
		}
	}
	return nil
}
