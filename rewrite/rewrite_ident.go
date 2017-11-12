package rewrite

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
)

// rewriteIdent converts TypeXXX to its replacement defined in typeMap.
func (s *Spec) rewriteIdent(pkg *Package) error {
	for _, node := range pkg.Files {
		var used []string
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				if x.Obj == nil || x.Obj.Kind != ast.Typ {
					return false
				}
				to, ok := s.TypeMap[x.Name]
				if !ok {
					return false
				}
				x.Name = to.Expr

				if len(to.Import) == 0 {
					return false
				}
				for _, im := range to.Import {
					var found bool
					for _, prev := range used {
						if im == prev {
							found = true
							break
						}
					}
					if !found {
						used = append(used, im)
					}
				}
				return false
			}
			return true
		})
		for _, im := range used {
			astutil.AddImport(pkg.FileSet, node, im)
		}
	}
	return nil
}
