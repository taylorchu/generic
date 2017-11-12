package rewrite

type Type struct {
	Expr   string
	Import []string
}

type Spec struct {
	TypeMap map[string]Type `yaml:"typeMap"`

	Name   string
	Import string
	Local  bool
}

type Config struct {
	Spec []*Spec
}

func (c *Config) RewritePackage() error {
	for _, s := range c.Spec {
		pkg, err := s.parse()
		if err != nil {
			return err
		}
		resetAST := func(pkg *Package) error {
			return pkg.Reset()
		}

		// Apply AST changes and refresh.
		for _, rewriteFunc := range []func(*Package) error{
			s.rewritePackageName,
			s.removePlaceholder,
			s.rewriteIdent,
			s.rewriteTopLevelIdent,
			resetAST,
			s.typeCheck,
			s.writePackage,
		} {
			err := rewriteFunc(pkg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
