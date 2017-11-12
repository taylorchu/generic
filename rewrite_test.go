package generic

import "testing"

func TestRewritePackage(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/test/pkg/basic", "result", map[string]Target{
		"Type": Target{Ident: "int64"},
	}, "test/output/basic")
}

func TestRewritePackageVendoring(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/test/pkg/basic", "result", map[string]Target{
		"Type": Target{Ident: "vendoring.Number", Import: "github.com/taylorchu/generic/test/pkg/vendoring"},
	}, "test/output/vendoring")
}

func TestRewritePackageMethod(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/test/pkg/method", "result", map[string]Target{
		"Type2": Target{Ident: "generic.Target", Import: "github.com/taylorchu/generic"},
	}, "test/output/method")
}

func TestRewritePackageInternal(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/test/pkg/basic", "internal/result", map[string]Target{
		"Type": Target{Ident: "int64"},
	}, "test/output/internal")
}

func TestRewritePackageDotRename(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/test/pkg/rename", ".result", map[string]Target{
		"Type": Target{Ident: "int64"},
	}, "test/output/dot_rename")
}

func TestRewritePackageQueue(t *testing.T) {
	testRewritePackage(t, "github.com/taylorchu/generic/test/pkg/queue", "result", map[string]Target{
		"Type":      Target{Ident: "int64"},
		"TypeQueue": Target{Ident: "FIFO"},
	}, "test/output/queue")
}

func TestRewritePackageDotQueue(t *testing.T) {
	testRewritePackageWithInput(t, "github.com/taylorchu/generic/test/pkg/queue", ".result", map[string]Target{
		"Type":      Target{Ident: "Data"},
		"TypeQueue": Target{Ident: "FIFO"},
	},
		"test/input/data",
		"test/output/dot_queue",
	)
}

func TestRewritePackageDotQueuePrefix(t *testing.T) {
	testRewritePackageWithInput(t, "github.com/taylorchu/generic/test/pkg/queue", ".result", map[string]Target{
		"Type": Target{Ident: "Data"},
	},
		"test/input/data",
		"test/output/dot_queue_prefix",
	)
}

func TestRewritePackageDotContainer(t *testing.T) {
	testRewritePackageWithInput(t, "github.com/taylorchu/generic/test/pkg/container", ".result", map[string]Target{
		"Type":          Target{Ident: "*Data"},
		"TypeContainer": Target{Ident: "Box"},
	},
		"test/input/data",
		"test/output/dot_container",
	)
}

func TestRewritePackageDotContainerUpdate(t *testing.T) {
	testRewritePackageWithInput(t, "github.com/taylorchu/generic/test/pkg/container", ".result", map[string]Target{
		"Type":          Target{Ident: "*Data"},
		"TypeContainer": Target{Ident: "Box"},
	},
		"test/input/container_updated",
		"test/output/dot_container",
	)
}

func TestRewritePackageDotRenameUnresolved(t *testing.T) {
	testRewritePackageWithInput(t, "github.com/taylorchu/generic/test/pkg/rename", ".result", map[string]Target{
		"Type": Target{Ident: "Data"},
	},
		"test/input/data_unresolved",
		"test/output/dot_rename_unresolved",
	)
}
