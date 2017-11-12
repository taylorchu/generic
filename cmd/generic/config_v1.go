package main

import (
	"strings"

	"github.com/taylorchu/generic/rewrite"
)

// ConfigV1 specifies rewrite inputs.
type ConfigV1 struct {
	FromPkgPath string
	PkgPath     string
	Local       bool
	TypeMap     map[string]Target
}

// NewConfig creates a new rewrite context.
func NewConfig(pkgPath, newPkgPath string, rules ...string) (*ConfigV1, error) {
	c := &ConfigV1{
		FromPkgPath: pkgPath,
	}
	if strings.HasPrefix(newPkgPath, ".") {
		c.Local = true
		c.PkgPath = strings.TrimPrefix(newPkgPath, ".")
	} else {
		c.PkgPath = newPkgPath
	}

	typeMap, err := ParseTypeMap(rules)
	if err != nil {
		return nil, err
	}
	c.TypeMap = typeMap

	return c, nil
}

func (c1 *ConfigV1) RewritePackage() error {
	spec := &rewrite.Spec{
		Name:    c1.PkgPath,
		Local:   c1.Local,
		Import:  c1.FromPkgPath,
		TypeMap: make(map[string]rewrite.Type),
	}
	for from, to := range c1.TypeMap {
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
