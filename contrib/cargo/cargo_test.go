package cargo

import (
	"testing"

	"github.com/R9295/schloss/core"

	"github.com/stretchr/testify/assert"
)

// TODO make packages realistic

func TestPoetryCollectPackage(t *testing.T) {
	pkgs := []LockfilePackage{
		{
			Name:         "parserkiosk",
			Version:      "0.3.0",
			Dependencies: []string{"deno_core", "deno_runtime"},
		},
		{
			Name:         "deno_core",
			Version:      "22.6.0",
			Dependencies: []string{"tokio"},
		},
	}
	collectedPkgs := collectPackagesAsMap(pkgs)
	assert.Equal(t, pkgs[0], collectedPkgs["parserkiosk"])
	assert.Equal(t, pkgs[1], collectedPkgs["deno_core"])
}

func TestDiffPackagesRemovePackage(t *testing.T) {
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{},
		},
		{
			Name:         "deno_core",
			Version:      "22.6.0",
			Dependencies: []string{},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{},
		},
	}}
	diffList := DiffLockfiles(&oldLockfile, &newLockfile)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.REMOVED,
		MetaType: core.DEPENDENCY,
		Name:     "deno_core",
		Text:     "",
	})
}

func TestDiffPackagesAddPackage(t *testing.T) {
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "deno_core",
			Version:      "42.0",
			Dependencies: []string{},
		},
		{
			Name:         "tokio",
			Version:      "22.6.0",
			Dependencies: []string{},
		},
	}}
	diffList := DiffLockfiles(&oldLockfile, &newLockfile)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.ADDED,
		MetaType: core.DEPENDENCY,
		Name:     "deno_core",
		Text:     "version=42.0",
	})
}


func TestDiffPackagesPackageVersion(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.0",
		Source: "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}

	newPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.1",
		Source: "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}
	diffList := make([]core.Diff, 0)
	diffList = diffPackages(&oldPkg, &newPkg, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Text:     "(old)version=0.3.0 & (new)version=0.3.1",
	})
}

func TestDiffPackagesPackageSouce(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.0",
		Source: "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}

	newPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.0",
		Source: "registry+https://githubb.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}
	diffList := make([]core.Diff, 0)
	diffList = diffPackages(&oldPkg, &newPkg, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Text:     "(old)source=registry+https://github.com/rust-lang/crates.io-index & (new)source=registry+https://githubb.com/rust-lang/crates.io-index",
	})
}

func TestDiffPackagesPackageChecksum(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.0",
		Source: "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}

	newPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.0",
		Source: "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}
	diffList := make([]core.Diff, 0)
	diffList = diffPackages(&oldPkg, &newPkg, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Text:     "(old)checksum=Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581 & (new)checksum=7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	})
}
