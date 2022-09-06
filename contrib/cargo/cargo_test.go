package cargo

import (
	"fmt"
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
	collectedPkgs, _ := collectPackages(pkgs)
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
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateRemovedDependencyDiff("deno_core"))
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
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateAddedDependencyDiff("deno_core", "42.0", "rootPkg"))
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
	assert.Equal(t, diffList[0],
		core.GenerateDependencyFieldDiff(
			"parserkiosk",
			"version",
			oldPkg.Version,
			newPkg.Version,
		))
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
	assert.Equal(t, diffList[0],
		core.GenerateDependencyFieldDiff(
			"parserkiosk",
			"source",
			oldPkg.Source,
			newPkg.Source,
		))
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
	assert.Equal(t, diffList[0], core.GenerateDependencyFieldDiff(
		"parserkiosk",
		"checksum",
		oldPkg.Checksum,
		newPkg.Checksum,
	))
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
	assert.Equal(t, diffList[0], core.GenerateAddedSubDependencyDiff(
		"new",
		"parserkiosk",
		"",
	))
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
	assert.Equal(t, diffList[0], core.GenerateRemovedSubDependencyDiff(
		"remove",
		"parserkiosk",
	))
}

func TestNoDuplicateModifiedSubDependencyWhenAdding(t *testing.T) {
	/*
		When adding a new pkg which has a shared, existing sub-dependency with an existing pkg
		if the sub-dependency is modified(eg. version bump), make sure the modification diff
		is only for the existing pkg and not for the added.
	*/
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.1",
			Dependencies: []string{},
		},
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{"sub_dep"},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.2",
			Dependencies: []string{},
		},
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{"sub_dep"},
		},
		{
			Name:         "deno_core",
			Version:      "42.0",
			Dependencies: []string{"sub_dep"},
		},
	},
	}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 3)
	assert.Equal(
		t, []core.Diff{
			core.GenerateDependencyFieldDiff("sub_dep", "version", "0.1", "0.2"),
			core.GenerateAddedDependencyDiff("deno_core", "42.0", "rootPkg"),
			core.GenerateModifiedSubDependencyDiff("sub_dep", "tokio"),
		},
		diffList)
}

func TestNoDuplicateModifiedSubDependencyWhenRemoving(t *testing.T) {
	/*
		When adding a new pkg which has a shared, existing sub-dependency with a removed pkg
		if the sub-dependency is modified(eg. version bump), make sure the modification diff
		is only for the existing pkg and not for the removed.
	*/
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.2",
			Dependencies: []string{},
		},
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{"sub_dep"},
		},
		{
			Name:         "deno_core",
			Version:      "42.0",
			Dependencies: []string{"sub_dep"},
		},
	},
	}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.1",
			Dependencies: []string{},
		},
		{
			Name:         "tokio",
			Version:      "42.0",
			Dependencies: []string{"sub_dep"},
		},
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	fmt.Println(diffList)
	assert.Equal(t, len(diffList), 3)
	assert.Equal(
		t, []core.Diff{
			core.GenerateDependencyFieldDiff("sub_dep", "version", "0.2", "0.1"),
			core.GenerateRemovedDependencyDiff("deno_core"),
			core.GenerateModifiedSubDependencyDiff("sub_dep", "tokio"),
		},
		diffList)
}
