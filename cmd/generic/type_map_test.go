package main

import (
	"reflect"
	"testing"

	"github.com/taylorchu/generic/rewrite"
)

func TestParseTypeMap(t *testing.T) {
	for _, test := range []struct {
		in   string
		want map[string]rewrite.Type
	}{
		{
			in: "",
		},
		{
			in: "T->V",
		},
		{
			in: "Type->",
		},
		{
			in: " Type  -> OtherType   ",
			want: map[string]rewrite.Type{
				"Type": {
					Expr: "OtherType",
				},
			},
		},
		{
			in: "Type->:OtherType",
		},
		{
			in: "Type->github.com/go:",
		},
		{
			in: "Type->  github.com/go :  go.OtherType ",
			want: map[string]rewrite.Type{
				"Type": {
					Import: []string{"github.com/go"},
					Expr:   "go.OtherType",
				},
			},
		},
	} {
		tm, err := ParseTypeMap([]string{test.in})
		if test.want == nil {
			if err == nil {
				t.Fatalf("expect error, got %v", tm)
			}
		} else {
			if !reflect.DeepEqual(tm, test.want) {
				t.Fatalf("expect %v, got %s", test.want, err)
			}
		}
	}
}
