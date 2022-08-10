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
	packages := collectPackagesAsMap(&lockfile)
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
	diffList = diffPackages(&oldPkg, &newPkg, diffList)
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
	diffList := DiffLockfiles(&oldLockfile, &newLockfile)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.REMOVED,
		MetaType: core.DEPENDENCY,
		Name:     "black",
		Parent:   "",
	})
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
	diffList := DiffLockfiles(&oldLockfile, &newLockfile)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.ADDED,
		MetaType: core.DEPENDENCY,
		Name:     "black",
		Parent:   "",
		Version:  "22.6.0",
	})
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
	diffList = diffPackages(&old, &new, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.REMOVED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "MarkupSafe",
		Parent:   "jinja2",
	})
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
	diffList = diffPackages(&old, &new, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.DependencyDiff{
		Type:     core.ADDED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "MarkupSafe",
		Parent:   "jinja2",
		Version:  ">=2.0",
	})
}

