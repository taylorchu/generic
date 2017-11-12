package rewrite

import (
	"errors"
	"os"
	"path/filepath"
)

// rewritePackageName sets current package name.
func (s *Spec) rewritePackageName(pkg *Package) error {
	pkgName := filepath.Base(s.Name)
	if s.Local {
		pkgName = os.Getenv("GOPACKAGE")
		if pkgName == "" {
			return errors.New("GOPACKAGE cannot be empty")
		}
	}
	for _, node := range pkg.Files {
		node.Name.Name = pkgName
	}
	return nil
}
