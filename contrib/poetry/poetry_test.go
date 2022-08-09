package poetry

import (
	"testing"

	"github.com/R9295/schloss/contrib/toml"
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
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     "parserkiosk",
		Text:     "(old)version=0.3.0 & (new)version=0.3.1",
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
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.REMOVED,
		MetaType: core.DEPENDENCY,
		Name:     "black",
		Text:     "",
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
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.ADDED,
		MetaType: core.DEPENDENCY,
		Name:     "black",
		Text:     "version=22.6.0",
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
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.REMOVED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "MarkupSafe",
		Text:     "of jinja2",
	})
}

func TestDiffPackagesPackageModifySubDependency(t *testing.T) {
	old := LockfilePackage{
		Name:    "jinja2",
		Version: "42.0",
		Dependencies: map[string]interface{}{
			"MarkupSafe": ">=2.0",
		},
	}
	new := LockfilePackage{
		Name:    "jinja2",
		Version: "42.0",
		Dependencies: map[string]interface{}{
			"MarkupSafe": ">=3.0",
		},
	}
	diffList := make([]core.Diff, 0)
	diffList = diffPackages(&old, &new, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "MarkupSafe",
		Text:     "of jinja2 | (old)version=>=2.0 & (new)version=>=3.0",
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
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.ADDED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "MarkupSafe",
		Text:     "of jinja2",
	})
}

func TestDiffPackagesPackageModifySubDependencyWithVersionMap(t *testing.T) {
	oldLockfileText := `[[package]]
	name = "django"
	version = "4.0.6"
	[package.dependencies]
	tzdata = {version = "*", markers = "sys_platform == \"win32\""}`
	var oldLockfile Lockfile
	toml.DecodeToml(oldLockfileText, &oldLockfile)
	new := LockfilePackage{
		Name:    "django",
		Version: "4.0.6",
		Dependencies: map[string]interface{}{
			"tzdata": ">=2.0",
		},
	}
	diffList := make([]core.Diff, 0)
	diffList = diffPackages(&oldLockfile.Package[0], &new, diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.SUB_DEPENDENCY,
		Name:     "tzdata",
		Text:     "of django | (old)version=* & (new)version=>=2.0",
	})
}
