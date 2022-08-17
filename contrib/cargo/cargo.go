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

type ParsedSubPackages map[string][]string

func collectPackages(lockFilePkgs []LockfilePackage) (map[string]LockfilePackage, map[string][]string) {
	packages := make(map[string]LockfilePackage)
	subPackages := make(map[string][]string)
	for _, pkg := range lockFilePkgs {
		packages[pkg.Name] = pkg
		for _, subPkg := range pkg.Dependencies {
			_, exists := subPackages[subPkg]
			if !exists {
				subPackages[subPkg] = []string{pkg.Name}
			} else {
				subPackages[subPkg] = append(subPackages[subPkg], pkg.Name)
			}
		}
	}
	return packages, subPackages
}

func diffSubPackages(oldSubPackages ParsedSubPackages, newSubPackages ParsedSubPackages) ParsedSubPackages {
	diff := make(ParsedSubPackages)
	for pkgName, oldValue := range oldSubPackages {
		newValue, exists := newSubPackages[pkgName]
		if exists {
			diff[pkgName] = []string{}
			for _, oldDep := range oldValue {
				if index := core.IndexOf(newValue, oldDep); index != -1 {
					diff[pkgName] = append(diff[pkgName], oldDep)
				}
			}
		}
	}
	return diff
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
	oldPkgs, oldSubPackages := collectPackages(oldLockfileToml.Package)
	newPkgs, newSubPackages := collectPackages(newLockfileToml.Package)
	diffedSubPackages := diffSubPackages(oldSubPackages, newSubPackages)
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
	var subPkgDiff []core.Diff

	for _, diff := range *diffList {
		if diff.GetType() == core.MODIFIED {
			depName := diff.GetName()
			parents, exists := diffedSubPackages[depName]
			if exists {
				for _, parent := range parents {
					subPkgDiff = append(subPkgDiff, core.GenerateModifiedSubDependencyDiff(depName, parent))
				}
			}
		}
	}
	*diffList = append(*diffList, subPkgDiff...)
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
