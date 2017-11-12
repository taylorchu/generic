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
			astutil.AddImport(pkg.FileSet, node, im)
		}
	}
	return nil
}
