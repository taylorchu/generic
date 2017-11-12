package rewrite

import (
	"fmt"
	"go/ast"
)

// rewriteTopLevelIdent adds a prefix to top-level identifiers and their uses.
//
// This prevents name conflicts when a package is rewritten to $PWD.
func (s *Spec) rewriteTopLevelIdent(pkg *Package) error {
	if !s.Local {
		return nil
	}

	prefixIdent := func(name string) string {
		if name == "_" {
			// skip unnamed
			return "_"
		}
		return lintName(fmt.Sprintf("%s_%s", s.Name, name))
	}

	declMap := make(map[interface{}]string)

	for _, node := range pkg.Files {
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
							if to, ok := s.TypeMap[obj.Name]; ok && spec.Name.Name == to.Expr {
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
	for _, node := range pkg.Files {
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
	return nil
}
