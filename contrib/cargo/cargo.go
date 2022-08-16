package cargo

import (
	"github.com/R9295/schloss/contrib/toml"
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
func doesSubdependencyExist(dependencyList []string, check string) bool {
	for _, dependency := range dependencyList {
		if dependency == check {
			return true
		}
	}
	return false
}

func diffPackages(oldPkg *LockfilePackage, newPkg *LockfilePackage, diffList *[]core.Diff) {
	if oldPkg.Version != newPkg.Version {
		*diffList = append(*diffList,
			core.GenerateDependencyFieldDiff(newPkg.Name,
				"version",
				oldPkg.Version,
				newPkg.Version),
		)
	}
	if oldPkg.Checksum != newPkg.Checksum {
		*diffList = append(*diffList,
			core.GenerateDependencyFieldDiff(newPkg.Name,
				"checksum",
				oldPkg.Checksum,
				newPkg.Checksum),
		)
	}
	if oldPkg.Source != newPkg.Source {
		*diffList = append(*diffList,
			core.GenerateDependencyFieldDiff(newPkg.Name,
				"source",
				oldPkg.Source,
				newPkg.Source),
		)
	}
	for _, dependency := range oldPkg.Dependencies {
		if doesSubdependencyExist(newPkg.Dependencies, dependency) == false {
			*diffList = append(*diffList,
				core.GenerateRemovedSubDependencyDiff(dependency, oldPkg.Name),
			)
		}
	}
	for _, dependency := range newPkg.Dependencies {
		if doesSubdependencyExist(oldPkg.Dependencies, dependency) == false {
			*diffList = append(*diffList,
				core.GenerateAddedSubDependencyDiff(dependency, newPkg.Name, ""),
			)
		}
	}
}

func DiffLockfiles(oldLockfileToml *Lockfile, newLockfileToml *Lockfile, diffList *[]core.Diff) {
	oldPkgs := collectPackagesAsMap(oldLockfileToml.Package)
	newPkgs := collectPackagesAsMap(newLockfileToml.Package)
	for oldPkgName, oldPkg := range oldPkgs {
		newPkg, exists := newPkgs[oldPkgName]
		if !exists {
			*diffList = append(*diffList, core.GenerateRemovedDependencyDiff(oldPkgName))
		} else {
			diffPackages(&oldPkg, &newPkg, diffList)
			delete(newPkgs, oldPkgName)
		}
	}
	for newPkgName, newPkgValue := range newPkgs {
		*diffList = append(*diffList, core.GenerateAddedDependencyDiff(newPkgName, newPkgValue.Version))
	}
}

func Diff(oldLockfile *string, newLockfile *string, diffList *[]core.Diff) error {
	oldLockfileToml, err := toml.ParseLockfile[Lockfile](*oldLockfile)
	if err != nil {
		return err
	}
	newLockfileToml, err := toml.ParseLockfile[Lockfile](*newLockfile)
	if err != nil {
		return err
	}
	DiffLockfiles(&oldLockfileToml, &newLockfileToml, diffList)
	return nil
}
