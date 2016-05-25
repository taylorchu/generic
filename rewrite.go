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
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// Target represents replacement output.
type Target struct {
	Ident  string
	Import string
}

// rewritePkgName sets current package name.
func rewritePkgName(node ast.Node, pkgName string) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			x.Name.Name = pkgName
		}
		return true
	})
}

// rewriteIdent converts TypeXXX to its replacement defined in typeMap.
func rewriteIdent(node ast.Node, typeMap map[string]Target, fset *token.FileSet) {
	var file *ast.File
	var used []string
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Obj != nil && x.Obj.Kind == ast.Typ {
				if to, ok := typeMap[x.Name]; ok {
					x.Name = to.Ident
					if to.Import != "" {
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
					}
				}
			}
		case *ast.File:
			file = x
		}
		return true
	})
	if file != nil {
		for _, im := range used {
			astutil.AddImport(fset, file, im)
		}
	}
}

// removeTypeDecl removes type declarations defined in typeMap.
func removeTypeDecl(node ast.Node, typeMap map[string]Target) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.File:
			for i := len(x.Decls) - 1; i >= 0; i-- {
				var remove bool
				genDecl, ok := x.Decls[i].(*ast.GenDecl)
				if !ok {
					continue
				}
				if genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if _, ok := typeMap[typeSpec.Name.Name]; ok {
						remove = true
					}
				}
				if remove {
					x.Decls = append(x.Decls[:i], x.Decls[i+1:]...)
				}
			}
		}
		return true
	})
}

// walkSource visits all .go files in a package path except tests.
func walkSource(pkgPath string, sourceFunc func(string) error) error {
	pkgPath = fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), pkgPath)
	return filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != pkgPath {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		return sourceFunc(path)
	})
}

// RewritePackage applies type replacements on a package in GOPATH, and saves results as a new package in $PWD.
//
// If there is a dir with the same name as newPkgPath, it will first be removed. It is possible to re-run this
// to update a generic package.
func RewritePackage(pkgPath string, newPkgPath string, typeMap map[string]Target) error {
	var err error
	if newPkgPath != "." {
		err = os.RemoveAll(newPkgPath)
		if err != nil {
			return err
		}

		err = os.MkdirAll(newPkgPath, 0777)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				os.RemoveAll(newPkgPath)
			}
		}()
	}

	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	err = walkSource(pkgPath, func(path string) error {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		files[filepath.Base(path)] = f
		return nil
	})
	if err != nil {
		return err
	}

	// Gather ast.File to create ast.Package.
	// ast.NewPackage will try to resolve unresolved identifiers.
	ast.NewPackage(fset, files, nil, nil)

	// Find out new package name.
	newPkgName := filepath.Base(newPkgPath)
	if newPkgPath == "." {
		gopackage := os.Getenv("GOPACKAGE")
		if gopackage == "" {
			return errors.New("GOPACKAGE cannot be empty")
		}
		newPkgName = gopackage
	}

	// Apply AST changes and refresh.
	buf := new(bytes.Buffer)
	for basename, f := range files {
		rewritePkgName(f, newPkgName)
		removeTypeDecl(f, typeMap)
		rewriteIdent(f, typeMap, fset)

		// AST in dirty state; refresh
		buf.Reset()
		err = printer.Fprint(buf, fset, f)
		if err != nil {
			return err
		}
		files[basename], err = parser.ParseFile(fset, "", buf, 0)
		if err != nil {
			return err
		}
	}

	conf := types.Config{Importer: importer.Default()}
	for basename, f := range files {
		// Type check.
		_, err = conf.Check("", fset, []*ast.File{f}, nil)
		if err != nil {
			return err
		}

		// Print ast to file.
		var dest *os.File
		dest, err = os.Create(fmt.Sprintf("%s/%s", newPkgPath, basename))
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
