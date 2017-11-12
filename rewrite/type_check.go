package rewrite

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/types"
	"strings"
)

func (s *Spec) typeCheck(pkg *Package) error {
	if s.Local {
		return nil
	}

	var allFiles []*ast.File
	for _, f := range pkg.Files {
		allFiles = append(allFiles, f)
	}
	allFileSets := pkg.FileSet

	var errType []error
	conf := types.Config{
		Importer: importer.For("source", nil),
		Error: func(err error) {
			// Ignore undeclared name error because we want developers to use this tool
			// during development process.
			if strings.HasPrefix(err.(types.Error).Msg, "undeclared name: ") {
				return
			}
			errType = append(errType, err)
		},
	}
	conf.Check("", allFileSets, allFiles, nil)
	if len(errType) > 0 {
		for _, err := range errType {
			fmt.Println(err)
		}
		return errType[0]
	}
	return nil
}
