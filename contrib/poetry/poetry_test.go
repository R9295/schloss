package poetry

import (
	"fmt"
	"strings"
	"testing"

	"github.com/R9295/schloss/core"
	"github.com/brianvoe/gofakeit/v6"

	"github.com/stretchr/testify/assert"
)

func getRandomName() string {
	return strings.Replace(
		fmt.Sprintf("%s-%s", gofakeit.HipsterWord(), gofakeit.Animal()),
		" ",
		"-",
		-1,
	)
}

func getRandomDependency() (string, string) {
	return getRandomName(), gofakeit.AppVersion()
}
func getRandLockfilePkg() LockfilePackage {
	pkg := LockfilePackage{
		Name:         getRandomName(),
		Version:      gofakeit.AppVersion(),
		Dependencies: map[string]interface{}{},
	}
	depsAmount := gofakeit.IntRange(0, 15)
	for i := 0; i < depsAmount; i++ {
		name, version := getRandomDependency()
		pkg.Dependencies[name] = version
	}
	return pkg
}

func TestPoetryCollectPackage(t *testing.T) {
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	lockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}
	packages, _ := collectPackages(&lockfile)
	assert.Equal(t, packages[pkgOne.Name], lockfile.Package[0])
	assert.Equal(t, packages[pkgTwo.Name], lockfile.Package[1])

}

func TestDiffPackagesPackageVersion(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Version = "666.666"
	diffList := make([]core.Diff, 0)
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeDependencyFieldDiff(
		oldPkg.Name, "version", oldPkg.Version, newPkg.Version,
	))
}

func TestDiffPackagesRemovePackage(t *testing.T) {
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	oldLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeRemovedDependencyDiff(pkgTwo.Name))
}

func TestDiffPackagesAddPackage(t *testing.T) {
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	oldLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
	}}
	newLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(
		t,
		diffList[0],
		core.MakeAddedDependencyDiff(pkgTwo.Name, pkgTwo.Version, "rootPkg"),
	)
}

func TestDiffPackagesPackageRemoveSubDependency(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := LockfilePackage{
		Name:         oldPkg.Name,
		Version:      oldPkg.Version,
		Dependencies: map[string]interface{}{},
	}
	for k, v := range oldPkg.Dependencies {
		newPkg.Dependencies[k] = v
	}
	oldPkg.Dependencies["to-remove"] = "0.1"
	diffList := make([]core.Diff, 0)
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeRemovedSubDependencyDiff("to-remove", oldPkg.Name))
}

func TestDiffPackagesPackageAddSubDependency(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := LockfilePackage{
		Name:         oldPkg.Name,
		Version:      oldPkg.Version,
		Dependencies: map[string]interface{}{},
	}
	for k, v := range oldPkg.Dependencies {
		newPkg.Dependencies[k] = v
	}
	newPkg.Dependencies["added"] = "0.1"
	diffList := make([]core.Diff, 0)
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeAddedSubDependencyDiff(
		"added",
		oldPkg.Name,
		"0.1",
	))
}

func TestNoDuplicateModifiedSubDependencyWhenAdding(t *testing.T) {
	pkgOne := getRandLockfilePkg()
	pkgOne.Version = "0.1"
	pkgTwo := getRandLockfilePkg()
	pkgTwo.Dependencies[pkgOne.Name] = "0.1"
	oldLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}
	pkgThree := getRandLockfilePkg()
	pkgThree.Dependencies[pkgOne.Name] = "0.1"
	pkgOne.Version = "0.2"
	newLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
		pkgThree,
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 3)
	assert.ElementsMatch(
		t, []core.Diff{
			core.MakeDependencyFieldDiff(pkgOne.Name, "version", "0.1", "0.2"),
			core.MakeAddedDependencyDiff(pkgThree.Name, pkgThree.Version, "rootPkg"),
			core.MakeModifiedSubDependencyDiff(pkgOne.Name, pkgTwo.Name),
		},
		diffList)
}

func TestNoDuplicateModifiedSubDependencyWhenRemoving(t *testing.T) {
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	pkgThree := getRandLockfilePkg()
	pkgTwo.Dependencies[pkgOne.Name] = "0.2"
	pkgThree.Dependencies[pkgOne.Name] = "0.2"
	pkgOne.Version = "0.2"
	oldLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
		pkgThree,
	}}
	pkgOne.Version = "0.1"
	newLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	assert.Equal(t, len(diffList), 3)
	assert.ElementsMatch(
		t, []core.Diff{
			core.MakeDependencyFieldDiff(pkgOne.Name, "version", "0.2", "0.1"),
			core.MakeRemovedDependencyDiff(pkgThree.Name),
			core.MakeModifiedSubDependencyDiff(pkgOne.Name, pkgTwo.Name),
		},
		diffList)
}
