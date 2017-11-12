// Package generic generates package with type replacements.
package generic

import (
	"github.com/taylorchu/generic/rewrite"
)

// RewritePackage applies type replacements on a package in $GOPATH, and saves results as a new package in $PWD.
//
// If there is a dir with the same name as newPkgPath, it will first be removed. It is possible to re-run this
// to update a generic package.
func RewritePackage(ctx *Context) error {
	spec := &rewrite.Spec{
		Name:    ctx.PkgPath,
		Package: ctx.PkgName,
		Local:   ctx.Local,
		Import:  ctx.FromPkgPath,
		TypeMap: make(map[string]rewrite.Type),
	}
	for from, to := range ctx.TypeMap {
		spec.TypeMap[from] = rewrite.Type{
			Ident:  to.Ident,
			Import: to.Import,
		}
	}
	c := rewrite.Config{
		Spec: []*rewrite.Spec{spec},
	}
	return c.RewritePackage()
}
