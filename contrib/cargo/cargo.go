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

func generateFieldDiff(pkgName string, fieldName string, oldVal string, newVal string) core.Diff {
	return core.Diff{
		Type:     core.MODIFIED,
		MetaType: core.DEPENDENCY,
		Name:     pkgName,
		Text:     fmt.Sprintf("(old)%s=%s & (new)%s=%s", fieldName, oldVal, fieldName, newVal),
	}
}
func diffPackages(oldPkg *LockfilePackage, newPkg *LockfilePackage, diffList []core.Diff) []core.Diff {
	if oldPkg.Version != newPkg.Version {
		diffList = append(diffList,
			generateFieldDiff(newPkg.Name,
					"version",
					oldPkg.Version,
					newPkg.Version),
			)
	}
	if oldPkg.Checksum != newPkg.Checksum {
		diffList = append(diffList,
			generateFieldDiff(newPkg.Name,
					"checksum",
					oldPkg.Checksum,
					newPkg.Checksum),
			)
	}
	if oldPkg.Source != newPkg.Source {
		diffList = append(diffList,
			generateFieldDiff(newPkg.Name,
					"source",
					oldPkg.Source,
					newPkg.Source),
			)
	}
	return diffList
}

func DiffLockfiles(oldLockfileToml *Lockfile, newLockfileToml *Lockfile) []core.Diff {
	diffList := make([]core.Diff, 0)
	oldPkgs := collectPackagesAsMap(oldLockfileToml.Package)
	newPkgs := collectPackagesAsMap(newLockfileToml.Package)
	for oldPkgName := range oldPkgs {
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
