package rewrite

import (
	"fmt"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
)

func (s *Spec) writePackage(pkg *Package) error {
	fset := token.NewFileSet()
	writeOutput := func() error {
		for path, f := range pkg.Files {
			path := filepath.Join(s.Name, filepath.Base(path))
			if s.Local {
				path = fmt.Sprintf("%s_%s", s.Name, filepath.Base(path))
			}
			// Print ast to file.
			dest, err := os.Create(path)
			if err != nil {
				return err
			}
			defer dest.Close()

			err = format.Node(dest, fset, f)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if s.Local {
		return writeOutput()
	}

	err := os.RemoveAll(s.Name)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.Name, 0777)
	if err != nil {
		return err
	}
	err = writeOutput()
	if err != nil {
		os.RemoveAll(s.Name)
		return err
	}

	return nil
}
