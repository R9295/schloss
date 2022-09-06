package poetry

import (
	"testing"

	"github.com/R9295/schloss/core"

	"github.com/stretchr/testify/assert"
)

func TestPoetryCollectPackage(t *testing.T) {
	lockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:    "parserkiosk",
			Version: "0.3.0",
			Dependencies: map[string]interface{}{
				"Jinja2":     "",
				"PyYAML":     "",
				"python-box": "",
				"yamale":     "",
			},
		},
		{
			Name:         "black",
			Version:      "22.6.0",
			Dependencies: map[string]interface{}{},
		},
	}}
	packages, _ := collectPackages(&lockfile)
	assert.Equal(t, packages["parserkiosk"], lockfile.Package[0])
}

func TestDiffPackagesPackageVersion(t *testing.T) {
	oldPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.0",
	}

	newPkg := LockfilePackage{
		Name:    "parserkiosk",
		Version: "0.3.1",
	}
	diffList := make([]core.Diff, 0)
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateDependencyFieldDiff(
		"parserkiosk", "version", oldPkg.Version, newPkg.Version,
	))
}

func TestDiffPackagesRemovePackage(t *testing.T) {
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "jinja2",
			Version:      "42.0",
			Dependencies: map[string]interface{}{},
		},
		{
			Name:         "black",
			Version:      "22.6.0",
			Dependencies: map[string]interface{}{},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "jinja2",
			Version:      "42.0",
			Dependencies: map[string]interface{}{},
		},
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateRemovedDependencyDiff("black"))
}

func TestDiffPackagesAddPackage(t *testing.T) {
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "jinja2",
			Version:      "42.0",
			Dependencies: map[string]interface{}{},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "jinja2",
			Version:      "42.0",
			Dependencies: map[string]interface{}{},
		},
		{
			Name:         "black",
			Version:      "22.6.0",
			Dependencies: map[string]interface{}{},
		},
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateAddedDependencyDiff("black", "22.6.0", "rootPkg"))
}

func TestDiffPackagesPackageRemoveSubDependency(t *testing.T) {
	old := LockfilePackage{
		Name:    "jinja2",
		Version: "42.0",
		Dependencies: map[string]interface{}{
			"MarkupSafe": ">=2.0",
		},
	}
	new := LockfilePackage{
		Name:         "jinja2",
		Version:      "42.0",
		Dependencies: map[string]interface{}{},
	}
	diffList := make([]core.Diff, 0)
	diffPackages(&old, &new, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateRemovedSubDependencyDiff("MarkupSafe", "jinja2"))
}

func TestDiffPackagesPackageAddSubDependency(t *testing.T) {
	old := LockfilePackage{
		Name:         "jinja2",
		Version:      "42.0",
		Dependencies: map[string]interface{}{},
	}
	new := LockfilePackage{
		Name:    "jinja2",
		Version: "42.0",
		Dependencies: map[string]interface{}{

			"MarkupSafe": ">=2.0",
		},
	}
	diffList := make([]core.Diff, 0)
	diffPackages(&old, &new, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.GenerateAddedSubDependencyDiff(
		"MarkupSafe",
		"jinja2",
		">=2.0",
	))
}

func TestNoDuplicateModifiedSubDependencyWhenAdding(t *testing.T) {
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.1",
			Dependencies: map[string]interface{}{},
		},
		{
			Name:    "jinja2",
			Version: "42.0",
			Dependencies: map[string]interface{}{
				"sub_dep": "",
			},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.2",
			Dependencies: map[string]interface{}{},
		},
		{
			Name:    "jinja2",
			Version: "42.0",
			Dependencies: map[string]interface{}{
				"sub_dep": "",
			},
		},
		{
			Name:    "black",
			Version: "22.6.0",
			Dependencies: map[string]interface{}{
				"sub_dep": "",
			},
		},
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 3)
	assert.Equal(
		t, []core.Diff{
			core.GenerateDependencyFieldDiff("sub_dep", "version", "0.1", "0.2"),
			core.GenerateAddedDependencyDiff("black", "22.6.0", "rootPkg"),
			core.GenerateModifiedSubDependencyDiff("sub_dep", "jinja2"),
		},
		diffList)
}

func TestNoDuplicateModifiedSubDependencyWhenRemoving(t *testing.T) {
	oldLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.2",
			Dependencies: map[string]interface{}{},
		},
		{
			Name:    "jinja2",
			Version: "42.0",
			Dependencies: map[string]interface{}{
				"sub_dep": "",
			},
		},
		{
			Name:    "black",
			Version: "22.6.0",
			Dependencies: map[string]interface{}{
				"sub_dep": "",
			},
		},
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		{
			Name:         "sub_dep",
			Version:      "0.1",
			Dependencies: map[string]interface{}{},
		},
		{
			Name:    "jinja2",
			Version: "42.0",
			Dependencies: map[string]interface{}{
				"sub_dep": "",
			},
		},
	}}

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 3)
	assert.Equal(
		t, []core.Diff{
			core.GenerateDependencyFieldDiff("sub_dep", "version", "0.2", "0.1"),
			core.GenerateRemovedDependencyDiff("black"),
			core.GenerateModifiedSubDependencyDiff("sub_dep", "jinja2"),
		},
		diffList)
}
