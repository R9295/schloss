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
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.REMOVED,
		MetaType: core.DEPENDENCY,
		Name:     "deno_core",
		Parent:   "",
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
			Version:      "42.0",
			Dependencies: []string{},
		},
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.ADDED,
		MetaType: core.DEPENDENCY,
		Name:     "deno_core",
		Parent:   "",
		Version:  "42.0",
	})
}

func TestDiffPackagesPackageVersion(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:     "parserkiosk",
		Version:  "0.3.0",
		Source:   "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}

	newPkg := LockfilePackage{
		Name:     "parserkiosk",
		Version:  "0.3.1",
		Source:   "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.FieldDiff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Field:    "version",
		OldValue: oldPkg.Version,
		NewValue: newPkg.Version,
	})
}

func TestDiffPackagesPackageSouce(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:     "parserkiosk",
		Version:  "0.3.0",
		Source:   "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}

	newPkg := LockfilePackage{
		Name:     "parserkiosk",
		Version:  "0.3.0",
		Source:   "registry+https://githubb.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.FieldDiff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Field:    "source",
		OldValue: oldPkg.Source,
		NewValue: newPkg.Source,
	})
}

func TestDiffPackagesPackageChecksum(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:     "parserkiosk",
		Version:  "0.3.0",
		Source:   "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}

	newPkg := LockfilePackage{
		Name:     "parserkiosk",
		Version:  "0.3.0",
		Source:   "registry+https://github.com/rust-lang/crates.io-index",
		Checksum: "7a8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
	}
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.FieldDiff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Field:    "checksum",
		NewValue: newPkg.Checksum,
		OldValue: oldPkg.Checksum,
	})
}

func TestDiffPackagesSubDependencyAdd(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:         "parserkiosk",
		Version:      "0.3.0",
		Source:       "registry+https://github.com/rust-lang/crates.io-index",
		Checksum:     "Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
		Dependencies: []string{"something"},
	}
	newPkg := LockfilePackage{
		Name:         "parserkiosk",
		Version:      "0.3.0",
		Source:       "registry+https://github.com/rust-lang/crates.io-index",
		Checksum:     "Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
		Dependencies: []string{"something", "new"},
	}
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.ADDED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "new",
		Parent:   "parserkiosk",
		Version:  "",
	})
}

func TestDiffPackagesSubDependencyRemove(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:         "parserkiosk",
		Version:      "0.3.0",
		Source:       "registry+https://github.com/rust-lang/crates.io-index",
		Checksum:     "Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
		Dependencies: []string{"something", "remove"},
	}
	newPkg := LockfilePackage{
		Name:         "parserkiosk",
		Version:      "0.3.0",
		Source:       "registry+https://github.com/rust-lang/crates.io-index",
		Checksum:     "Aa8325f63a7d4774dd041e363b2409ed1c5cbbd0f867795e661df066b2b0a581",
		Dependencies: []string{"something"},
	}
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.REMOVED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "remove",
		Parent:   "parserkiosk",
	})
}
