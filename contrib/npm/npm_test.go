package npm

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
		Version:        gofakeit.AppVersion(),
		Resolved:       gofakeit.URL(),
		Integrity:      gofakeit.Regex("[a-zA-Z0-9]{64}"),
		Dependencies:   map[string]string{},
		ParentPackages: []string{},
	}
	/* depsAmount := gofakeit.IntRange(0, 15)
	for i := 0; i < depsAmount; i++ {
		pkg.Dependencies = append(pkg.Dependencies, getRandomName())
	} */
	return pkg
}

// when parsing package name remove prefix node_modules/
func getRandomLockfile() Lockfile {
	lockfile := Lockfile{
		Name:            getRandomName(),
		Version:         gofakeit.AppVersion(),
		LockfileVersion: gofakeit.AppVersion(),
		Packages:        map[string]LockfilePackage{},
	}
	for i := 0; i < 3; {
		lockfile.Packages[getRandomName()] = getRandLockfilePkg()
		i += 1
	}
	return lockfile
}

func deepCopyPackage(p LockfilePackage) LockfilePackage {
	copyDependencies := make(map[string]string)

	for k, v := range p.Dependencies {
		copyDependencies[k] = v
	}

	copyParentPackages := []string{}
	copy(copyParentPackages, p.ParentPackages)

	return LockfilePackage{
		Version:        p.Version,
		Integrity:      p.Integrity,
		Dependencies:   copyDependencies,
		Resolved:       p.Resolved,
		ParentPackages: copyParentPackages,
	}
}

func deepCopyLockfile(lockfile Lockfile) Lockfile {
	copyLockfile := Lockfile{
		Name:            lockfile.Name,
		Version:         lockfile.Version,
		LockfileVersion: lockfile.LockfileVersion,
		Packages:        map[string]LockfilePackage{},
	}
	for k, v := range lockfile.Packages {
		copyLockfile.Packages[k] = deepCopyPackage(v)
	}
	return copyLockfile
}

func TestDiffLockfilesWithoutChanges(t *testing.T) {

	oldLockfile := getRandomLockfile()
	newLockfile := oldLockfile
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	assert.Equal(t, len(diffList), 0)
}

func TestDiffLockfilesModifiedVersion(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)
	testPackage := getRandLockfilePkg()
	testPackage.Version = "6.2.3"
	packageName := getRandomName()
	oldLockfile.Packages[packageName] = testPackage
	modifiedTestPackage := testPackage
	modifiedTestPackage.Version = "7.2.3"
	newLockfile.Packages[packageName] = modifiedTestPackage
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeDependencyFieldDiff(packageName, "version", "6.2.3", "7.2.3")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffLockfilesModifiedResolved(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)
	testPackage := getRandLockfilePkg()
	originalResolved := testPackage.Resolved
	packageName := getRandomName()
	oldLockfile.Packages[packageName] = testPackage
	modifiedTestPackage := testPackage
	modifiedTestPackage.Resolved = originalResolved + "t"
	newLockfile.Packages[packageName] = modifiedTestPackage
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeDependencyFieldDiff(packageName, "resolved", originalResolved, originalResolved+"t")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffLockfilesModifiedIntegrity(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)
	testPackage := getRandLockfilePkg()
	originalIntegrity := testPackage.Integrity
	packageName := getRandomName()
	oldLockfile.Packages[packageName] = testPackage
	modifiedTestPackage := testPackage
	modifiedTestPackage.Integrity = originalIntegrity + "T"
	newLockfile.Packages[packageName] = modifiedTestPackage
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeDependencyFieldDiff(packageName, "integrity", originalIntegrity, originalIntegrity+"T")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffLockfilesAddPackage(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)
	addedPackage := getRandLockfilePkg()
	packageName := getRandomName()
	newLockfile.Packages[packageName] = addedPackage
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeAddedDependencyDiff(packageName, addedPackage.Version, newLockfile.Name)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffLockfilesRemovePackage(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)
	addedPackage := getRandLockfilePkg()
	packageName := getRandomName()
	oldLockfile.Packages[packageName] = addedPackage
	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeRemovedDependencyDiff(packageName)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffLockfilesAddSubDependency(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)

	parentPackage := getRandLockfilePkg()
	parentPackageName := getRandomName()

	oldLockfile.Packages[parentPackageName] = parentPackage
	newLockfile.Packages[parentPackageName] = deepCopyPackage(parentPackage)
	subDependencyName := getRandomName()

	newLockfile.Packages[parentPackageName].Dependencies[subDependencyName] = "2.0.0"

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeAddedSubDependencyDiff(subDependencyName, parentPackageName, "2.0.0")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffLockfilesRemovSubDependency(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)

	parentPackage := getRandLockfilePkg()
	parentPackageName := getRandomName()

	oldLockfile.Packages[parentPackageName] = parentPackage
	newLockfile.Packages[parentPackageName] = deepCopyPackage(parentPackage)
	subDependencyName := getRandomName()

	oldLockfile.Packages[parentPackageName].Dependencies[subDependencyName] = "2.0.0"

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeRemovedSubDependencyDiff(subDependencyName, parentPackageName)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffPackagesAbsentFieldVersion(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)

	testPackage := getRandLockfilePkg()
	packageName := getRandomName()

	oldLockfile.Packages[packageName] = testPackage
	modifiedTestPackage := deepCopyPackage(testPackage)
	modifiedTestPackage.Version = ""
	newLockfile.Packages[packageName] = modifiedTestPackage

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeAbsentFieldDiff(packageName, "version")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}
func TestDiffPackagesAbsentFieldIntegrity(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)

	testPackage := getRandLockfilePkg()
	packageName := getRandomName()

	oldLockfile.Packages[packageName] = testPackage
	modifiedTestPackage := deepCopyPackage(testPackage)
	modifiedTestPackage.Integrity = ""
	newLockfile.Packages[packageName] = modifiedTestPackage

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeAbsentFieldDiff(packageName, "integrity")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}
func TestDiffPackagesAbsentFieldResolved(t *testing.T) {
	oldLockfile := getRandomLockfile()
	newLockfile := deepCopyLockfile(oldLockfile)

	testPackage := getRandLockfilePkg()
	packageName := getRandomName()

	oldLockfile.Packages[packageName] = testPackage
	modifiedTestPackage := deepCopyPackage(testPackage)
	modifiedTestPackage.Resolved = ""
	newLockfile.Packages[packageName] = modifiedTestPackage

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeAbsentFieldDiff(packageName, "resolved")
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffMetadataVersion(t *testing.T) {
	oldLockfile := getRandomLockfile()
	oldLockfile.Version = "2.2.1"
	newLockfile := deepCopyLockfile(oldLockfile)
	newLockfile.Version = "2.2.2"

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeModifiedMetadataDiff("version", oldLockfile.Version, newLockfile.Version)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffMetadataLockfileVersion(t *testing.T) {
	oldLockfile := getRandomLockfile()
	oldLockfile.LockfileVersion = "2.2.1"
	newLockfile := deepCopyLockfile(oldLockfile)
	newLockfile.LockfileVersion = "2.2.2"

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeModifiedMetadataDiff("lockfileVersion", oldLockfile.LockfileVersion, newLockfile.LockfileVersion)
	assert.Equal(t, len(diffList), 1)
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestDiffMetadataName(t *testing.T) {
	oldLockfile := getRandomLockfile()
	oldLockfile.Name = "my-project"
	newLockfile := deepCopyLockfile(oldLockfile)
	newLockfile.Name = "my-other-project"

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeModifiedMetadataDiff("name", oldLockfile.Name, newLockfile.Name)
	assert.Equal(t, 1, len(diffList))
	assert.Equal(t, expectedDiff, diffList[0])
}

func TestCollectPackages(t *testing.T) {
	/*
		test collectPackages on lockfile with nested dependency structure
		package 1 + 5 as parent package
		packages 2 + 3 subpackages of package 1
		pacakge 4 subpackage of package 3 and 5
	*/
	testLockfile := getRandomLockfile()

	package1 := getRandLockfilePkg()
	package2 := getRandLockfilePkg()
	package3 := getRandLockfilePkg()
	package4 := getRandLockfilePkg()
	package5 := getRandLockfilePkg()
	name1 := getRandomName()
	name2 := getRandomName()
	name3 := getRandomName()
	name4 := getRandomName()
	name5 := getRandomName()

	testLockfile.Packages[name1] = package1
	testLockfile.Packages[name2] = package2
	testLockfile.Packages[name3] = package3
	testLockfile.Packages[name4] = package4
	testLockfile.Packages[name5] = package5

	testLockfile.Packages[name1].Dependencies[name2] = "2.0.2"
	testLockfile.Packages[name1].Dependencies[name3] = "2.0.3"
	testLockfile.Packages[name3].Dependencies[name4] = "2.0.4"
	testLockfile.Packages[name5].Dependencies[name4] = "2.0.4"

	expectedLockfile := deepCopyLockfile(testLockfile)
	if entry, exists := expectedLockfile.Packages[name2]; exists {
		entry.ParentPackages = append(expectedLockfile.Packages[name2].ParentPackages, name1)
	} else {
		t.Fail()
	}
	if entry, exists := expectedLockfile.Packages[name3]; exists {
		entry.ParentPackages = append(expectedLockfile.Packages[name3].ParentPackages, name1)
	} else {
		t.Fail()
	}
	if entry, exists := expectedLockfile.Packages[name4]; exists {
		entry.ParentPackages = append(expectedLockfile.Packages[name4].ParentPackages, name3, name5)
	} else {
		t.Fail()
	}

	collectPackages(&testLockfile)
	assert.Equal(t, expectedLockfile, testLockfile)

}

func TestDiffLockfilesModifiedSubdependency(t *testing.T) {
	oldLockfile := getRandomLockfile()
	parentPackage := getRandLockfilePkg()
	dependency := getRandLockfilePkg()
	parentPackageName := getRandomName()
	dependencyName := getRandomName()
	oldLockfile.Packages[parentPackageName] = parentPackage
	oldLockfile.Packages[dependencyName] = dependency
	oldLockfile.Packages[parentPackageName].Dependencies[dependencyName] = "2.0.0"

	newLockfile := deepCopyLockfile(oldLockfile)
	newLockfile.Packages[parentPackageName].Dependencies[dependencyName] = "2.0.1"

	var diffList []core.Diff
	DiffLockfiles(&oldLockfile, &newLockfile, &diffList)
	expectedDiff := core.MakeModifiedSubDependencyDiff(dependencyName, parentPackageName)
	assert.Equal(t, 1, len(diffList))
	assert.Equal(t, expectedDiff, diffList[0])
}

/*

(x) create functions to create random diffs

(x) unit test on small single methods => reak up any of the methods? Unit tests => for collect packages check if parent dependency field is correct

(x) lockfile name, version, lockfile version

(x) MakeModifiedSubDependencyDiff

Bei MakeModifiedSubDependencyDiff gibt es nur versions unterschiede, values werden nicht an core weitergegeben und tauchen nicht in der ausgabe auf

(x) happy path without changes
field changes
(x) modify field in dependency (for each field) [version, resolved, integrity]
(x) remove field in dependency (for each field) [version, resolved, integrity]
( ) add field in dependency (for each field) [version, resolved, integrity]
( ) what about double field? How would json get parsed?

(x) add entire new dependency
(x) remove entire dependency

(x) add subdependency
(x) remove subdependency

( )NEW TASK => do it in caro
	Refactor AddedDependencyDiff to AddedPackage and print parent packages so that it is clear wether it is a dependency or subdependency

( ) special case: added dependency is already there => dependency duplicate (what if 2 packages need different version of package?)

added subdependency => add package and add package name do existing pacakge dependencies list, 2 options:
	subdependency is listed in package-lock
	subdependency is not listed in package-lock, random name has been added


added subdependency means adding a entry in package dependencies list



What if supdependency package has been removed but dependency entry still in dependencies list(s)

Definition Subdependency => Dependency that is listed in an other dependency dependencies list

Definition Dependency => Dependency that does not appear in an other dependency dependencies list + appears under packages ""



When an entire package has been deleted/ added/ modified we call it Change of Dependency even if that Dependency is technically a Subdependency
=> No, only Dependencies that show up in package.json should be dependencies so dependencies that show up under
pacakges{
	""{
		dependencies: {
			asdfasfasdf
			asdfasfasdf
			asdfasdfasd
		}
	}
}
They have only one parent package which is the root package all others are subdependencies

What does added Subdependency mean? Added in dependencies list of a package or added

Why can I remove a dependency that is still listed as subdependency in an other package and clean install works? => Output if dependency is
in someone elses depencencies list would be useful

check for changes in packages ""

=> treat change of dependencies list als FieldDiff
=> added pacakges are added dependencies or subdependencies depending if the have parent packages that are not equal to root pacakge


context is key, some changes are not bad but in combination with others they might be and you can't tell from one look
which other dependencies are affected by the removal/ modification of one, easy to overlook something



special case:
	added dependency that is treated as subdependency but is not a actual subdependency basically a dependency that
	is not being used by anything => big warning, check for unused dependency

*/