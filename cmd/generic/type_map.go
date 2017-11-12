package main

import (
	"errors"
	"strings"

	"github.com/taylorchu/generic/rewrite"
)

// ParseTypeMap parses raw strings to type replacements.
func ParseTypeMap(args []string) (map[string]rewrite.Type, error) {
	typeMap := make(map[string]rewrite.Type)

	for _, arg := range args {
		part := strings.Split(arg, "->")

		if len(part) != 2 {
			return nil, errors.New("RULE must be in form of `TypeXXX->OtherType`")
		}

		var (
			from = strings.TrimSpace(part[0])
			to   = strings.TrimSpace(part[1])
		)

		if !strings.HasPrefix(from, "Type") {
			return nil, errors.New("REPL type must start with `Type`")
		}

		var t rewrite.Type
		if strings.Contains(to, ":") {
			toPart := strings.Split(to, ":")

			if len(toPart) != 2 {
				return nil, errors.New("REPL type must be in form of DESTPATH:OtherType")
			}

			t.Import = strings.TrimSpace(toPart[0])
			t.Ident = strings.TrimSpace(toPart[1])
			if strings.Count(t.Ident, ".") != 1 {
				return nil, errors.New("REPL type must contain one `.`")
			}
		} else {
			t.Ident = to
			if strings.Count(t.Ident, ".") != 0 {
				return nil, errors.New("REPL type must not contain `.`")
			}
		}
		if t.Ident == "" {
			return nil, errors.New("REPL type cannot be empty")
		}

		typeMap[from] = t
	}
	return typeMap, nil
}
