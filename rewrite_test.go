package generic

import "testing"

func TestRewritePackage(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/fixture/basic", "result", map[string]Target{
		"Type": Target{Ident: "int64"},
	}, "output/basic")
}

func TestRewritePackageMethod(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/fixture/method", "result", map[string]Target{
		"Type2": Target{Ident: "generic.Target", Import: "github.com/taylorchu/generic"},
	}, "output/method")
}

func TestRewritePackageInternal(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/fixture/basic", "internal/result", map[string]Target{
		"Type": Target{Ident: "int64"},
	}, "output/internal")
}

func TestRewritePackageDot(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/fixture/rename", ".result", map[string]Target{
		"Type": Target{Ident: "int64"},
	}, "output/rename")
}
