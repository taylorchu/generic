package rewrite

import "testing"

func TestRewritePackage(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/basic",
			TypeMap: map[string]Type{
				"Type": Type{Expr: "int64"},
			},
		},
	}}
	testRewritePackage(t, c, "test/output/basic")
}

func TestRewritePackageVendoring(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/basic",
			TypeMap: map[string]Type{
				"Type": Type{Expr: "vendoring.Number", Import: []string{"github.com/taylorchu/generic/rewrite/test/pkg/vendoring"}},
			},
		},
	}}
	testRewritePackage(t, c, "test/output/vendoring")
}

func TestRewritePackageMethod(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/method",
			TypeMap: map[string]Type{
				"Type2": Type{Expr: "vendoring.Number", Import: []string{"github.com/taylorchu/generic/rewrite/test/pkg/vendoring"}},
			},
		},
	}}
	testRewritePackage(t, c, "test/output/method")
}

func TestRewritePackageInternal(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "internal/result",
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/basic",
			TypeMap: map[string]Type{
				"Type": Type{Expr: "int64"},
			},
		},
	}}
	testRewritePackage(t, c, "test/output/internal")
}

func TestRewritePackageRenameLocal(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Local:  true,
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/rename",
			TypeMap: map[string]Type{
				"Type": Type{Expr: "int64"},
			},
		},
	}}
	testRewritePackage(t, c, "test/output/rename_local")
}

func TestRewritePackageQueue(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/queue",
			TypeMap: map[string]Type{
				"Type":      Type{Expr: "int64"},
				"TypeQueue": Type{Expr: "FIFO"},
			},
		},
	}}
	testRewritePackage(t, c, "test/output/queue")
}

func TestRewritePackageQueueLocal(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Local:  true,
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/queue",
			TypeMap: map[string]Type{
				"Type":      Type{Expr: "Data"},
				"TypeQueue": Type{Expr: "FIFO"},
			},
		},
	}}
	testRewritePackageWithInput(t, c, "test/input/data", "test/output/queue_local")
}

func TestRewritePackageQueuePrefixLocal(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Local:  true,
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/queue",
			TypeMap: map[string]Type{
				"Type": Type{Expr: "Data"},
			},
		},
	}}
	testRewritePackageWithInput(t, c, "test/input/data", "test/output/queue_prefix_local")
}

func TestRewritePackageContainerLocal(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Local:  true,
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/container",
			TypeMap: map[string]Type{
				"Type":          Type{Expr: "*Data"},
				"TypeContainer": Type{Expr: "Box"},
			},
		},
	}}
	testRewritePackageWithInput(t, c, "test/input/data", "test/output/container_local")
}

func TestRewritePackageContainerLocalUpdate(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Local:  true,
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/container",
			TypeMap: map[string]Type{
				"Type":          Type{Expr: "*Data"},
				"TypeContainer": Type{Expr: "Box"},
			},
		},
	}}
	testRewritePackageWithInput(t, c, "test/input/container_updated", "test/output/container_local")
}

func TestRewritePackageRenameUnresolvedLocal(t *testing.T) {
	c := &Config{Spec: []*Spec{
		{
			Name:   "result",
			Local:  true,
			Import: "github.com/taylorchu/generic/rewrite/test/pkg/rename",
			TypeMap: map[string]Type{
				"Type": Type{Expr: "Data"},
			},
		},
	}}
	testRewritePackageWithInput(t, c, "test/input/data_unresolved", "test/output/rename_unresolved_local")
}
