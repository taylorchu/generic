package rewrite

import (
	"go/ast"
)

// removePlaceholder removes type declarations defined in typeMap.
func (s *Spec) removePlaceholder(pkg *Package) error {
	declMap := make(map[interface{}]struct{})
	for _, node := range pkg.Files {
		for i := len(node.Decls) - 1; i >= 0; i-- {
			var remove bool
			switch decl := node.Decls[i].(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						_, ok := s.TypeMap[spec.Name.Name]
						if !ok {
							continue
						}
						_, ok = spec.Type.(*ast.Ident)
						if !ok {
							continue
						}
						remove = true
						declMap[spec] = struct{}{}
					}
				}
			}
			if remove {
				node.Decls = append(node.Decls[:i], node.Decls[i+1:]...)
			}
		}
	}
	// If a type placeholder is removed, its linked methods should be removed too.
	// This works like go interface because now the replaced types need to implement these methods.
	for _, node := range pkg.Files {
		for i := len(node.Decls) - 1; i >= 0; i-- {
			var remove bool
			switch decl := node.Decls[i].(type) {
			case *ast.FuncDecl:
				if decl.Recv == nil {
					continue
				}
				var obj *ast.Object
				switch expr := decl.Recv.List[0].Type.(type) {
				case *ast.StarExpr:
					obj = expr.X.(*ast.Ident).Obj
				case *ast.Ident:
					obj = expr.Obj
				}
				if obj == nil || obj.Decl == nil {
					continue
				}
				_, ok := declMap[obj.Decl]
				if !ok {
					continue
				}
				remove = true
			}
			if remove {
				node.Decls = append(node.Decls[:i], node.Decls[i+1:]...)
			}
		}
	}
	return nil
}
