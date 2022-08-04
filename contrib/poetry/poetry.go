package poetry

import (
	"fmt"
	"reflect"

	"github.com/R9295/schloss/contrib/toml"
	"github.com/R9295/schloss/core"
)

type LockfileMetaTypedataFile struct {
	File string
	Hash string
}

type LockfilePackage struct {
	Name         string
	Version      string
	Dependencies map[string]interface{}
}

type Lockfile struct {
	Package      []LockfilePackage
	MetaTypedata struct {
		PythonVersions string `toml:"python-versions"`
		ContentHash    string `toml:"content-hash"`
		LockVersion    string `toml:"lock-version"`
		Files          map[string][]LockfileMetaTypedataFile
	}
}

func collectPackagesAsMap(lockFileToml *Lockfile) map[string]LockfilePackage {
	packages := make(map[string]LockfilePackage)
	for _, pkg := range lockFileToml.Package {
		packages[pkg.Name] = pkg
	}
	return packages
}

func extractVersionValue(version interface{}) string {
	value := reflect.ValueOf(version)
	if value.Kind() == reflect.String {
		return value.String()
	}
	for _, key := range value.MapKeys() {
		if key.String() == "version" {
			return fmt.Sprintf("%s", value.MapIndex(key))
		}
	}
	return ""
}
func diffPackages(oldPkg *LockfilePackage, newPkg *LockfilePackage, diffList []core.Diff) []core.Diff {
	if oldPkg.Version != newPkg.Version {
		diffList = append(diffList,
			core.Diff{
				Type:     core.MODIFIED,
				MetaType: core.DEPENDENCY,
				Name:     newPkg.Name,
				Text: fmt.Sprintf("(old)version=%s & (new)version=%s",
					oldPkg.Version,
					newPkg.Version),
			})
	}
	for oldPkgDep, oldPkgDepVersion := range oldPkg.Dependencies {
		newPkgDepVersion, exists := newPkg.Dependencies[oldPkgDep]
		if !exists {
			diffList = append(diffList,
				core.Diff{
					Type:     core.REMOVED,
					MetaType: core.SUB_DEPENDENCY,
					Name:     oldPkgDep,
					Text:     fmt.Sprintf("of %s", oldPkg.Name),
				})
		} else {
			newPkgDepVersionValue := extractVersionValue(newPkgDepVersion)
			oldPkgDepVersionValue := extractVersionValue(oldPkgDepVersion)
			if oldPkgDepVersionValue != newPkgDepVersionValue {
				diffList = append(diffList,
					core.Diff{
						Type:     core.MODIFIED,
						MetaType: core.SUB_DEPENDENCY,
						Name:     oldPkgDep,
						Text: fmt.Sprintf("of %s | (old)version=%s & (new)version=%s",
							newPkg.Name,
							oldPkgDepVersionValue,
							newPkgDepVersionValue),
					})
				}
			delete(newPkg.Dependencies, oldPkgDep)
		}
	}
	for newPkgDep, newPkgDepVersion := range newPkg.Dependencies {
		diffList = append(diffList,
			core.Diff{
				Type:     core.ADDED,
				MetaType: core.SUB_DEPENDENCY,
				Name:     newPkgDep,
				Text: fmt.Sprintf("of %s | version=%s",
					newPkg.Name,
					extractVersionValue(newPkgDepVersion)),
			})
	}
	return diffList
}

func DiffLockfiles(oldLockfileToml *Lockfile, newLockfileToml *Lockfile) []core.Diff {
	diffList := make([]core.Diff, 0)
	oldPkgs := collectPackagesAsMap(oldLockfileToml)
	newPkgs := collectPackagesAsMap(newLockfileToml)
	for oldPkgName, oldPkgValue := range oldPkgs {
		newPkgValue, exists := newPkgs[oldPkgName]
		if !exists {
			diffList = append(diffList, core.Diff{
				Type:     core.REMOVED,
				MetaType: core.DEPENDENCY,
				Name:     oldPkgName,
			})
		} else {
			diffList = diffPackages(&oldPkgValue, &newPkgValue, diffList)
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

func ParseLockfiles(oldLockfile string, newLockfile string) (Lockfile, Lockfile) {
	var newLockfileToml Lockfile
	var oldLockfileToml Lockfile
	toml.DecodeToml(newLockfile, &newLockfileToml)
	toml.DecodeToml(oldLockfile, &oldLockfileToml)
	return oldLockfileToml, newLockfileToml
}
