package poetry

import (
	"fmt"
	"reflect"

	"github.com/R9295/schloss/contrib/toml"
	"github.com/R9295/schloss/core"
)

type RootTool struct {
	Name string
}
type RootFile struct {
	Tool map[string]RootTool
}
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

type ParsedSubPackages map[string][]string

func collectPackages(lockFile *Lockfile) (map[string]LockfilePackage, ParsedSubPackages) {
	packages := make(map[string]LockfilePackage)
	subPackages := make(ParsedSubPackages)
	for _, pkg := range lockFile.Package {
		packages[pkg.Name] = pkg
		for subPkg := range pkg.Dependencies {
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

func diffSubPackages(
	oldSubPackages ParsedSubPackages,
	newSubPackages ParsedSubPackages,
) ParsedSubPackages {
	diff := make(ParsedSubPackages)
	for pkgName, oldValue := range oldSubPackages {
		newValue, exists := newSubPackages[pkgName]
		if exists {
			diff[pkgName] = make([]string, 0)
			for _, oldDep := range oldValue {
				if index := core.IndexOf(newValue, oldDep); index != -1 {
					diff[pkgName] = append(diff[pkgName], oldDep)
				}
			}
		}
	}
	return diff
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
func diffPackages(oldPkg *LockfilePackage, newPkg *LockfilePackage, diffList *[]core.Diff) {
	if oldPkg.Version != newPkg.Version {
		*diffList = append(
			*diffList,
			core.GenerateDependencyFieldDiff(
				newPkg.Name,
				"version",
				oldPkg.Version,
				newPkg.Version,
			),
		)
	}
	for oldPkgDep := range oldPkg.Dependencies {
		_, exists := newPkg.Dependencies[oldPkgDep]
		if !exists {
			*diffList = append(*diffList,
				core.GenerateRemovedSubDependencyDiff(oldPkgDep, oldPkg.Name),
			)
		} else {
			delete(newPkg.Dependencies, oldPkgDep)
		}
	}
	for newPkgDep, newPkgDepVersion := range newPkg.Dependencies {
		*diffList = append(*diffList,
			core.GenerateAddedSubDependencyDiff(newPkgDep,
				newPkg.Name,
				extractVersionValue(newPkgDepVersion)),
		)
	}
}

func DiffLockfiles(
	oldLockfileToml *Lockfile,
	newLockfileToml *Lockfile,
	diffList *[]core.Diff,
	rootPkg string,
) {
	oldPkgs, oldSubPackages := collectPackages(oldLockfileToml)
	newPkgs, newSubPackages := collectPackages(newLockfileToml)
	diffedSubPackages := diffSubPackages(oldSubPackages, newSubPackages)
	for oldPkgName, oldPkgValue := range oldPkgs {
		newPkgValue, exists := newPkgs[oldPkgName]
		if !exists {
			*diffList = append(*diffList, core.GenerateRemovedDependencyDiff(oldPkgName))
		} else {
			diffPackages(&oldPkgValue, &newPkgValue, diffList)
			delete(newPkgs, oldPkgName)
		}
	}
	for newPkgName, newPkgValue := range newPkgs {
		*diffList = append(
			*diffList,
			core.GenerateAddedDependencyDiff(newPkgName, newPkgValue.Version, rootPkg),
		)
	}

	var subPkgDiff []core.Diff
	for _, diff := range *diffList {
		if diff.GetType() == core.MODIFIED {
			depName := diff.GetName()
			parents, exists := diffedSubPackages[depName]
			if exists {
				for _, parent := range parents {
					subPkgDiff = append(
						subPkgDiff,
						core.GenerateModifiedSubDependencyDiff(depName, parent),
					)
				}
			}
		}
	}
	*diffList = append(*diffList, subPkgDiff...)
}

func GetRootPackageName(rootFile *string) (string, error) {
	rootToml, err := toml.ParseLockfile[RootFile](*rootFile)
	if err != nil {
		return "", err
	}
	poetryRoot, exists := rootToml.Tool["poetry"]
	if exists == false {
		return "", fmt.Errorf("root file did not contain tool.poetry")
	}
	if poetryRoot.Name == "" {
		return "", fmt.Errorf("root file did not declare tool.poetry.name")
	}
	return rootToml.Tool["poetry"].Name, nil
}
func Diff(rootFile *string, oldLockfile *string, newLockfile *string, diffList *[]core.Diff) error {
	oldLockfileToml, err := toml.ParseLockfile[Lockfile](*oldLockfile)
	if err != nil {
		return err
	}
	newLockfileToml, err := toml.ParseLockfile[Lockfile](*newLockfile)
	if err != nil {
		return err
	}
	rootPkg, err := GetRootPackageName(rootFile)
	if err != nil {
		return err
	}
	DiffLockfiles(&oldLockfileToml, &newLockfileToml, diffList, rootPkg)
	return nil
}
