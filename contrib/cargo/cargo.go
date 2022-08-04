package cargo

import (
	"fmt"

	"github.com/R9295/schloss/core"
)

type LockfilePackage struct {
	Name         string
	Version      string
	Source       string
	Checksum     string
	Dependencies []string
}

type Lockfile struct {
	Package []LockfilePackage
}

func collectPackagesAsMap(lockFilePkgs []LockfilePackage) map[string]LockfilePackage {
	packages := make(map[string]LockfilePackage)
	for _, pkg := range lockFilePkgs {
		packages[pkg.Name] = pkg
	}
	return packages
}

func DiffLockfiles(oldLockfileToml *Lockfile, newLockfileToml *Lockfile) []core.Diff {
	diffList := make([]core.Diff, 0)
	oldPkgs := collectPackagesAsMap(oldLockfileToml.Package)
	newPkgs := collectPackagesAsMap(newLockfileToml.Package)
	for oldPkgName, _ := range oldPkgs {
		_, exists := newPkgs[oldPkgName]
		if !exists {
			diffList = append(diffList, core.Diff{
				Type:     core.REMOVED,
				MetaType: core.DEPENDENCY,
				Name:     oldPkgName,
			})
		} else {
			delete(newPkgs, oldPkgName)
		}
	}
	for newPkgName, newPkgValue := range newPkgs {
		diffList = append(diffList,
			core.Diff{
				Type:     core.ADDED,
				MetaType: core.DEPENDENCY,
				Name:     newPkgName,
				Text:     fmt.Sprintf("version=%s", newPkgValue.Version),
			})
	}
	return diffList
}
