package rewrite

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

type Package struct {
	Files   map[string]*ast.File
	FileSet *token.FileSet
}

func (p *Package) Reset() error {
	p.FileSet = token.NewFileSet()
	buf := new(bytes.Buffer)
	for name, f := range p.Files {
		buf.Reset()
		err := format.Node(buf, p.FileSet, f)
		if err != nil {
			return err
		}
		parsed, err := parser.ParseFile(p.FileSet, "", buf, 0)
		if err != nil {
			format.Node(os.Stderr, p.FileSet, f)
			return err
		}
		p.Files[name] = parsed
	}

	// Gather ast.File to create ast.Package.
	// ast.NewPackage will try to resolve unresolved identifiers.
	//
	// It will return errors because the importer is not provided.
	ast.NewPackage(p.FileSet, p.Files, nil, nil)
	return nil
}