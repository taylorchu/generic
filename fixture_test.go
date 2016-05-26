package generic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func testRewritePackage(t *testing.T, pkgPath, newPkgPath string, typeMap map[string]Target, expect string) {
	const dirname = "rewrite_test"
	err := os.MkdirAll(dirname, 0777)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirname)

	err = os.Chdir(dirname)
	if err != nil {
		t.Fatal(err)
	}

	if strings.HasPrefix(newPkgPath, ".") {
		err = os.Setenv("GOPACKAGE", "GOPACKAGE")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Unsetenv("GOPACKAGE")
	}

	err = RewritePackage(pkgPath, newPkgPath, typeMap)
	if err != nil {
		t.Fatal(err)
	}

	os.Chdir("..")
	assertEqualDir(t, dirname, expect)
}

func assertEqualDir(t *testing.T, path1, path2 string) {
	t.Log(path1, path2)
	fi1, err := ioutil.ReadDir(path1)
	if err != nil {
		t.Fatal(err)
	}
	fi2, err := ioutil.ReadDir(path2)
	if err != nil {
		t.Fatal(err)
	}
	if len(fi1) != len(fi2) {
		t.Fatalf("%s: %d, %s: %d", path1, len(fi1), path2, len(fi2))
	}

	for _, info := range fi1 {
		p1 := fmt.Sprintf("%s/%s", path1, info.Name())
		p2 := fmt.Sprintf("%s/%s", path2, info.Name())
		if info.IsDir() {
			assertEqualDir(t, p1, p2)
		} else {
			b1, err := ioutil.ReadFile(p1)
			if err != nil {
				t.Fatal(err)
			}
			b2, err := ioutil.ReadFile(p2)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b1, b2) {
				t.Fatalf("%s:\n%s, %s:\n%s", p1, b1, p2, b2)
			}
		}
	}
}
