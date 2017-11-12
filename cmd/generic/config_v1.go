package main

import (
	"strings"

	"github.com/taylorchu/generic/rewrite"
)

// ConfigV1 specifies rewrite inputs.
type ConfigV1 rewrite.Spec

// NewConfig creates a new rewrite context.
func NewConfig(pkgPath, newPkgPath string, rules ...string) (*ConfigV1, error) {
	c := &ConfigV1{
		Import: pkgPath,
	}
	if strings.HasPrefix(newPkgPath, ".") {
		c.Local = true
		c.Name = strings.TrimPrefix(newPkgPath, ".")
	} else {
		c.Name = newPkgPath
	}

	typeMap, err := ParseTypeMap(rules)
	if err != nil {
		return nil, err
	}
	c.TypeMap = typeMap

	return c, nil
}

func (c1 *ConfigV1) RewritePackage() error {
	spec := rewrite.Spec(*c1)
	c := rewrite.Config{
		Spec: []*rewrite.Spec{&spec},
	}
	return c.RewritePackage()
}
