package cargo

import (
	"testing"

	"github.com/R9295/schloss/core"

	"github.com/stretchr/testify/assert"
)

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
