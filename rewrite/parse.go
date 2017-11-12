package rewrite

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
)

func (s *Spec) parse() (*Package, error) {
	// NOTE: this package that we try to rewrite from should not contain vendor/.
	buildP, err := build.Import(s.Import, "", 0)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	files := make(map[string]*ast.File)
	for _, file := range buildP.GoFiles {
		path := filepath.Join(buildP.Dir, file)
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, err
		}
		files[path] = f
	}
	ast.NewPackage(fset, files, nil, nil)
	return &Package{
		Files:   files,
		FileSet: fset,
	}, nil
}
