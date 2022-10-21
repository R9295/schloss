package npm

import (
	"encoding/json"

	"github.com/R9295/schloss/core"
)

type Lockfile struct {
	Name            string                     `json:"name"`
	Version         string                     `json:"version"`
	LockfileVersion int                        `json:"lockfileVersion"`
	Packages        map[string]LockfilePackage `json:"packages"`
	SubPackages     map[string]LockfilePackage `json:"dependencies"`
}

type LockfilePackage struct {
	Version        string            `json:"version"`
	Resolved       string            `json:"resolved"`
	Integrity      string            `json:"integrity"`
	Dependencies   map[string]string `json:"dependencies"`
	ParentPackages []string
}

type ParsedSubPackages map[string][]string

func parseJson(lockfile string) (Lockfile, error) {
	data := Lockfile{}
	err := json.Unmarshal([]byte(lockfile), &data)
	return data, err
}

func collectPackages(lockFile *Lockfile, diffList *[]core.Diff) {
	for pkgName, pkg := range lockFile.Packages {
		for _, subPkg := range pkg.Dependencies {
			pgkSubdependency, exists := lockFile.SubPackages[subPkg]
			// find subdependencies in json object dependencies
			if exists {
				pgkSubdependency.ParentPackages = append(pgkSubdependency.ParentPackages, pkgName)
				lockFile.SubPackages[subPkg] = pgkSubdependency
			} else {
				*diffList = append(*diffList, core.MakeRemovedSubDependencyDiff(subPkg, pkgName))
			}
		}
	}

}

func diffPackages(oldPkgName string, oldPkg LockfilePackage, newPkg LockfilePackage, diffList *[]core.Diff) {
	// check for changed fields
	if oldPkg.Version != newPkg.Version {
		*diffList = append(*diffList, core.MakeDependencyFieldDiff(oldPkgName, "version", oldPkg.Version, newPkg.Version))
	}
	if oldPkg.Integrity != newPkg.Integrity {
		*diffList = append(*diffList, core.MakeDependencyFieldDiff(oldPkgName, "integrity", oldPkg.Integrity, newPkg.Integrity))
	}
	if oldPkg.Resolved != newPkg.Resolved {
		*diffList = append(*diffList, core.MakeDependencyFieldDiff(oldPkgName, "resolved", oldPkg.Resolved, newPkg.Resolved))
	}
	// check if subdendency List matches, quick check if list is empty
	for _, oldPkgDep := range oldPkg.Dependencies {
		exists := false
		for _, newPkgDep := range newPkg.Dependencies {
			if newPkgDep == oldPkgDep {
				exists = true
			}
		}
		if !exists {
			*diffList = append(*diffList, core.MakeRemovedSubDependencyDiff(oldPkgDep, oldPkgName))
		}
	}
}

func DiffLockfiles(
	oldLockfile *Lockfile,
	newLockfile *Lockfile,
	diffList *[]core.Diff,
) {
	collectPackages(newLockfile, diffList)
	collectPackages(oldLockfile, diffList)
	// check differences in packages
	for oldPkgName, oldPkg := range oldLockfile.Packages {
		newPkg, exists := newLockfile.Packages[oldPkgName]
		if !exists {
			*diffList = append(*diffList, core.MakeRemovedDependencyDiff(oldPkgName))
		} else {
			diffPackages(oldPkgName, oldPkg, newPkg, diffList)
			delete(newLockfile.Packages, oldPkgName)
		}
	}
	for newPackageName, newPkg := range newLockfile.Packages {
		*diffList = append(*diffList, core.MakeAddedDependencyDiff(newPackageName, newPkg.Version, newLockfile.Name))
	}
	// check differences in SubPackages
	for oldSubPkgName, oldSubPkg := range oldLockfile.SubPackages {
		newSubPkg, exists := newLockfile.SubPackages[oldSubPkgName]
		if !exists {
			// for every parent make a sub dependency diff => double diff entry because of diffPackages()?
			for _, parentPkgName := range oldSubPkg.ParentPackages {
				*diffList = append(*diffList, core.MakeRemovedSubDependencyDiff(oldSubPkgName, parentPkgName))
			}
		} else {
			diffPackages(oldSubPkgName, oldSubPkg, newSubPkg, diffList)
			delete(newLockfile.Packages, oldSubPkgName)
		}
	}
	for newSubPkgName, newSubPkg := range newLockfile.SubPackages {
		// for every parent
		for _, parentPkg := range newSubPkg.ParentPackages {
			*diffList = append(*diffList, core.MakeAddedSubDependencyDiff(newSubPkgName, parentPkg, newSubPkg.Version))
		}

	}
}

func Diff(rootFile *string, oldLockfile *string, newLockfile *string, diffList *[]core.Diff) error {
	oldLockfileJson, err := parseJson(*oldLockfile)
	if err != nil {
		return err
	}
	newLockfileJson, err := parseJson(*newLockfile)
	if err != nil {
		return err
	}
	DiffLockfiles(&oldLockfileJson, &newLockfileJson, diffList)
	return nil
}
