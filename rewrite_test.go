package generic

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestRewritePackage(t *testing.T) {
	defer os.RemoveAll("rewrite_test")

	err := RewritePackage("github.com/taylorchu/generic/test", "rewrite_test", map[string]Target{
		"Type2": Target{Ident: "generic.Target", Import: "github.com/taylorchu/generic"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRewritePackageDot(t *testing.T) {
	const dir = "rewrite_dot_test"

	defer os.RemoveAll(dir)
	os.Mkdir(dir, 0777)

	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir("..")

	os.Setenv("GOPACKAGE", dir)

	err = RewritePackage("github.com/taylorchu/generic/test", ".", map[string]Target{
		"Type2": Target{Ident: "generic.Target", Import: "github.com/taylorchu/generic"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWalkSource(t *testing.T) {
	want := []string{
		"test.go",
	}
	var got []string
	err := walkSource("github.com/taylorchu/generic/test", func(path string) error {
		got = append(got, filepath.Base(path))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("expect %v, got %v", want, got)
	}
}

func TestRewritePkgName(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{
			in: `package p`,
			want: `package p2
`,
		},
	} {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", test.in, 0)
		if err != nil {
			t.Fatal(err)
		}

		rewritePkgName(f, "p2")

		buf := new(bytes.Buffer)
		printer.Fprint(buf, fset, f)

		got := buf.String()
		if got != test.want {
			t.Fatalf("expect %s, got %s", test.want, got)
		}
	}
}

func TestRemoveTypeDecl(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{
			in: `package p

type Type int
`,
			want: `package p
`,
		},
	} {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", test.in, 0)
		if err != nil {
			t.Fatal(err)
		}

		removeTypeDecl(f, map[string]Target{"Type": {Ident: "OtherType"}})

		buf := new(bytes.Buffer)
		printer.Fprint(buf, fset, f)

		got := buf.String()
		if got != test.want {
			t.Fatalf("expect %s, got %s", test.want, got)
		}
	}
}

func TestRewriteIdent(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{
			in: `package p

func (Type Type) Type(Type Type) {

}
`,
			want: `package p

func (Type OtherType) Type(Type OtherType) {

}
`,
		},
		{
			in: `package p

type Struct struct {
	Type Type
}
`,
			want: `package p

type Struct struct {
	Type OtherType
}
`,
		},
	} {
		fset := token.NewFileSet()
		f1, err := parser.ParseFile(fset, "", `package p; type Type int64`, 0)
		if err != nil {
			t.Fatal(err)
		}
		f2, err := parser.ParseFile(fset, "", test.in, 0)
		if err != nil {
			t.Fatal(err)
		}

		ast.NewPackage(fset, map[string]*ast.File{
			"f1": f1,
			"f2": f2,
		}, nil, nil)

		rewriteIdent(f2, map[string]Target{
			"Type": {Ident: "OtherType"},
		}, fset)

		buf := new(bytes.Buffer)
		printer.Fprint(buf, fset, f2)

		got := buf.String()
		if got != test.want {
			t.Fatalf("expect %s, got %s", test.want, got)
		}
	}
}

func TestRewriteIdentWithImport(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{
			in: `package p

func add(_ Type, _ Type2) {}

func add2(_ Type2, _ Type) {}
`,
			want: `package p

import (
	"github.com/golang/test"
	"github.com/golang/test2"
)

func add(_ test.OtherType, _ test.OtherType2)	{}

func add2(_ test.OtherType2, _ test.OtherType)	{}
`,
		},
	} {
		fset := token.NewFileSet()
		f1, err := parser.ParseFile(fset, "", `package p; type Type int64; type Type2 int64`, 0)
		if err != nil {
			t.Fatal(err)
		}
		f2, err := parser.ParseFile(fset, "", test.in, 0)
		if err != nil {
			t.Fatal(err)
		}

		ast.NewPackage(fset, map[string]*ast.File{
			"f1": f1,
			"f2": f2,
		}, nil, nil)

		rewriteIdent(f2, map[string]Target{
			"Type":  {Ident: "test.OtherType", Import: "github.com/golang/test"},
			"Type2": {Ident: "test.OtherType2", Import: "github.com/golang/test2"},
		}, fset)

		buf := new(bytes.Buffer)
		printer.Fprint(buf, fset, f2)

		got := buf.String()
		if got != test.want {
			t.Fatalf("expect %s, got %s", test.want, got)
		}
	}
}
