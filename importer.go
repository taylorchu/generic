package generic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

type im struct {
	cache map[string]*types.Package
}

// NewImporter creates a new types.Importer.
//
// See https://github.com/golang/go/issues/11415.
// Many applications use the gcimporter package to read type information from compiled object files.
// There's no guarantee that those files are even remotely recent.
func NewImporter() types.Importer {
	return &im{
		cache: make(map[string]*types.Package),
	}
}

func (i *im) Import(pkgPath string) (*types.Package, error) {
	if pkgPath == "unsafe" {
		return types.Unsafe, nil
	}

	if pkg, ok := i.cache[pkgPath]; ok {
		return pkg, nil
	}

	fset := token.NewFileSet()
	var files []*ast.File
	err := walkSource(pkgPath, func(path string) error {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		files = append(files, f)
		return nil
	})
	if err != nil {
		return nil, err
	}

	conf := types.Config{
		Importer: i,
	}

	pkg, err := conf.Check(pkgPath, fset, files, nil)
	if err != nil {
		return nil, err
	}

	i.cache[pkgPath] = pkg
	return pkg, nil
}
