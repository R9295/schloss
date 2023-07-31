package npm

import (
	"encoding/json"
	"strconv"

	"github.com/R9295/schloss/core"
)

type Lockfile struct {
	Name            string                     `json:"name"`
	Version         string                     `json:"version"`
	LockfileVersion int                        `json:"lockfileVersion"`
	Packages        map[string]LockfilePackage `json:"packages"`
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

func collectPackages(lockFile *Lockfile) {
	for pkgName, pkg := range lockFile.Packages {
		for _, depName := range pkg.Dependencies {
			subPkg, exists := lockFile.Packages[depName]
			if exists {
				subPkg.ParentPackages = append(subPkg.ParentPackages, pkgName)
			}
		}
	}

}

func diffPackageFields(
	oldPkgName string,
	oldPkg LockfilePackage,
	newPkg LockfilePackage,
	diffList *[]core.Diff,
) {
	// check for changed fields
	if oldPkg.Version != newPkg.Version {
		if newPkg.Version == "" {
			*diffList = append(
				*diffList,
				core.MakeAbsentFieldDiff(oldPkgName, "version"),
			)
		} else {
			*diffList = append(
				*diffList,
				core.MakeDependencyFieldDiff(oldPkgName, "version", oldPkg.Version, newPkg.Version),
			)
		}

	}
	if oldPkg.Integrity != newPkg.Integrity {
		if newPkg.Integrity == "" {
			*diffList = append(
				*diffList,
				core.MakeAbsentFieldDiff(oldPkgName, "integrity"),
			)
		} else {
			*diffList = append(
				*diffList,
				core.MakeDependencyFieldDiff(
					oldPkgName,
					"integrity",
					oldPkg.Integrity,
					newPkg.Integrity,
				),
			)
		}
	}
	if oldPkg.Resolved != newPkg.Resolved {
		if newPkg.Resolved == "" {
			*diffList = append(
				*diffList,
				core.MakeAbsentFieldDiff(oldPkgName, "resolved"),
			)
		} else {
			*diffList = append(
				*diffList,
				core.MakeDependencyFieldDiff(oldPkgName, "resolved", oldPkg.Resolved, newPkg.Resolved),
			)
		}

	}
}

func diffPackageSubDependencies(
	oldPkgName string,
	oldPkg LockfilePackage,
	newPkg LockfilePackage,
	diffList *[]core.Diff,
) {

	for oldPkgDep, oldPkgDepVersion := range oldPkg.Dependencies {
		exists := false
		for newPkgDep, newPkgDepVersion := range newPkg.Dependencies {
			if newPkgDep == oldPkgDep {
				exists = true
				delete(newPkg.Dependencies, newPkgDep)
				if oldPkgDepVersion != newPkgDepVersion {
					*diffList = append(*diffList, core.MakeModifiedSubDependencyDiff(newPkgDep, oldPkgName))
				}
			}
		}
		if !exists {
			*diffList = append(*diffList, core.MakeRemovedSubDependencyDiff(oldPkgDep, oldPkgName))
		}
	}
	for newPkgDep, newPkgVersion := range newPkg.Dependencies {
		*diffList = append(*diffList, core.MakeAddedSubDependencyDiff(newPkgDep, oldPkgName, newPkgVersion))
	}
}

func diffMetadata(
	oldLockfile *Lockfile,
	newLockfile *Lockfile,
	diffList *[]core.Diff,
) {
	if oldLockfile.Version != newLockfile.Version {
		*diffList = append(*diffList, core.MakeModifiedMetadataDiff("version", oldLockfile.Version, newLockfile.Version))
	}
	if oldLockfile.Name != newLockfile.Name {
		*diffList = append(*diffList, core.MakeModifiedMetadataDiff("name", oldLockfile.Name, newLockfile.Name))
	}
	if oldLockfile.LockfileVersion != newLockfile.LockfileVersion {
		*diffList = append(*diffList, core.MakeModifiedMetadataDiff("lockfileVersion", strconv.Itoa(oldLockfile.LockfileVersion), strconv.Itoa(newLockfile.LockfileVersion)))
	}
}

func DiffLockfiles(
	oldLockfile *Lockfile,
	newLockfile *Lockfile,
	diffList *[]core.Diff,
) {
	collectPackages(newLockfile)
	collectPackages(oldLockfile)

	diffMetadata(oldLockfile, newLockfile, diffList)

	// check differences in packages
	for oldPkgName, oldPkg := range oldLockfile.Packages {
		newPkg, exists := newLockfile.Packages[oldPkgName]
		if !exists {
			*diffList = append(*diffList, core.MakeRemovedDependencyDiff(oldPkgName))
		} else {
			diffPackageFields(oldPkgName, oldPkg, newPkg, diffList)
			diffPackageSubDependencies(oldPkgName, oldPkg, newPkg, diffList)
			delete(newLockfile.Packages, oldPkgName)
		}
	}
	for newPackageName, newPkg := range newLockfile.Packages {
		*diffList = append(*diffList, core.MakeAddedDependencyDiff(newPackageName, newPkg.Version, newLockfile.Name))
		/* if sliceContains(newPkg.ParentPackages, newLockfile.Name) && len(newPkg.ParentPackages) == 1 {
			*diffList = append(*diffList, core.MakeAddedDependencyDiff(newPackageName, newPkg.Version, newLockfile.Name))
		} else {
			*diffList = append(*diffList, core.MakeAddedSubDependencyDiff(newPackageName, newPkg.Version, newLockfile.Name))
		} */

	}
}

/* func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
} */

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

/* func removeStringOfSlice(slice []string, removeString string) []string {
	newSlice := []string{}
	for _, element := range slice {
		if element != removeString {
			newSlice = append(newSlice, element)
		}
	}
	return newSlice
} */
