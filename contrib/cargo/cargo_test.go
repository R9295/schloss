package cargo

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

func getRandLockfilePkg() LockfilePackage {
	pkg := LockfilePackage{
		Name:         getRandomName(),
		Version:      gofakeit.AppVersion(),
		Checksum:     gofakeit.Regex("[a-zA-Z0-9]{64}"),
		Source:       gofakeit.URL(),
		Dependencies: []string{},
	}
	depsAmount := gofakeit.IntRange(0, 15)
	for i := 0; i < depsAmount; i++ {
		pkg.Dependencies = append(pkg.Dependencies, getRandomName())
	}
	return pkg
}

func TestPoetryCollectPackage(t *testing.T) {
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	pkgs := []LockfilePackage{
		pkgOne,
		pkgTwo,
	}
	collectedPkgs, _ := collectPackages(pkgs)
	assert.Equal(t, pkgs[0], collectedPkgs[pkgOne.Name])
	assert.Equal(t, pkgs[1], collectedPkgs[pkgTwo.Name])
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

func TestDiffPackagesPackageVersion(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Version = "666.666"
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0],
		core.MakeDependencyFieldDiff(
			oldPkg.Name,
			"version",
			oldPkg.Version,
			newPkg.Version,
		))
}

func TestDiffPackagesPackageSouce(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Source = gofakeit.URL()
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0],
		core.MakeDependencyFieldDiff(
			oldPkg.Name,
			"source",
			oldPkg.Source,
			newPkg.Source,
		))
}

func TestDiffPackagesPackageChecksum(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Checksum = "NEW_CHECKSUM"
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeDependencyFieldDiff(
		oldPkg.Name,
		"checksum",
		oldPkg.Checksum,
		newPkg.Checksum,
	))
}

func TestDiffPackagesSubDependencyAdd(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Dependencies = append(newPkg.Dependencies, "new-sub-depdendency")
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeAddedSubDependencyDiff(
		"new-sub-depdendency",
		oldPkg.Name,
		"",
	))
}

func TestDiffPackagesSubDependencyRemove(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	oldPkg.Dependencies = append(oldPkg.Dependencies, "to-remove")
	newPkg := oldPkg
	newPkg.Dependencies = newPkg.Dependencies[:len(newPkg.Dependencies)-1]
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0], core.MakeRemovedSubDependencyDiff(
		oldPkg.Dependencies[len(oldPkg.Dependencies)-1],
		oldPkg.Name,
	))
}

func TestNoDuplicateModifiedSubDependencyWhenAdding(t *testing.T) {
	/*
		When adding a new pkg which has a shared sub-dependency with an existing pkg
		if the sub-dependency is modified(eg. version bump), make sure the modification diff
		is only for the existing pkg and not for the added.
	*/
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	pkgOne.Version = "0.1"
	pkgTwo.Dependencies = append(pkgTwo.Dependencies, pkgOne.Name)
	oldLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}
	pkgOne.Version = "0.2"
	pkgThree := getRandLockfilePkg()
	pkgThree.Dependencies = append(pkgThree.Dependencies, pkgOne.Name)
	newLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
		pkgThree,
	},
	}
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
	/*
		When adding a new pkg which has a shared sub-dependency with a removed pkg
		if the sub-dependency is modified(eg. version bump), make sure the modification diff
		is only for the existing pkg and not for the removed.
	*/
	pkgOne := getRandLockfilePkg()
	pkgTwo := getRandLockfilePkg()
	pkgThree := getRandLockfilePkg()
	pkgTwo.Dependencies = append(pkgTwo.Dependencies, pkgOne.Name)
	pkgThree.Dependencies = append(pkgTwo.Dependencies, pkgOne.Name)
	pkgOne.Version = "0.2"
	oldLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
		pkgThree,
	},
	}
	pkgOne.Version = "0.1"
	newLockfile := Lockfile{Package: []LockfilePackage{
		pkgOne,
		pkgTwo,
	}}
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList, "rootPkg")
	fmt.Println(diffList)
	assert.Equal(t, len(diffList), 3)
	assert.ElementsMatch(
		t, []core.Diff{
			core.MakeDependencyFieldDiff(pkgOne.Name, "version", "0.2", "0.1"),
			core.MakeRemovedDependencyDiff(pkgThree.Name),
			core.MakeModifiedSubDependencyDiff(pkgOne.Name, pkgTwo.Name),
		},
		diffList)
}

func TestDiffPackagesAbsentFieldChecksum(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Checksum = ""
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0],
		core.MakeAbsentFieldDiff(
			oldPkg.Name,
			"checksum",
		))
}

func TestDiffPackagesAbsentFieldSource(t *testing.T) {
	oldPkg := getRandLockfilePkg()
	newPkg := oldPkg
	newPkg.Source = ""
	var diffList []core.Diff
	diffPackages(&oldPkg, &newPkg, &diffList)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, diffList[0],
		core.MakeAbsentFieldDiff(
			oldPkg.Name,
			"source",
		))
}
