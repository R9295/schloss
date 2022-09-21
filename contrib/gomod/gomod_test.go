package gomod

import (
	"fmt"
	"strings"
	"testing"

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

func getRandLockfilePkgString() string {
	pkg := ""
	repo := fmt.Sprintf(
		"%s.%s/%s/%s",
		gofakeit.DomainName(),
		gofakeit.DomainSuffix(),
		getRandomName(),
		getRandomName(),
	)
	version := gofakeit.AppVersion()
	pkg += fmt.Sprintf("%s v%s h1:%s\n", repo, version, "hash")
	pkg += fmt.Sprintf("%s v%s/go.mod h1:%s\n", repo, version, "hash")
	return pkg
}
func getRandLockfileString() (string, int) {
	file := ""
	depsAmount := gofakeit.IntRange(1, 15)
	for i := 0; i < depsAmount; i++ {
		file += getRandLockfilePkgString()
	}
	return file, depsAmount
}

func TestParseLockfile(t *testing.T) {
	lockfileString := `github.com/R9295/schloss v0.1 h1:hash
github.com/R9295/schloss v0.1/go.mod h1:hashModFile`
	parsedLockfile, err := ParseLockfile(&lockfileString)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, parsedLockfile, Lockfile{
		Packages: []LockfilePackage{
			{
				Name:            "github.com/R9295/schloss",
				Version:         "0.1",
				Checksum:        "h1:hash",
				ModfileChecksum: "h1:hashModFile",
			},
		},
	})
}

func TestParseLockfileMultiplePkgs(t *testing.T) {
	lockfileString, count := getRandLockfileString()
	fmt.Println(lockfileString)
	parsedLockfile, err := ParseLockfile(&lockfileString)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, len(parsedLockfile.Packages), count)
}
func TestParseLockfileIgnoreNewline(t *testing.T) {
	lockfileString := `


github.com/R9295/schloss v0.1 h1:hash
github.com/R9295/schloss v0.1/go.mod h1:hashModFile`
	parsedLockfile, err := ParseLockfile(&lockfileString)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, parsedLockfile, Lockfile{
		Packages: []LockfilePackage{
			{
				Name:            "github.com/R9295/schloss",
				Version:         "0.1",
				Checksum:        "h1:hash",
				ModfileChecksum: "h1:hashModFile",
			},
		},
	})
}

func TestParseLockfileBadFirstLine(t *testing.T) {
	lockfileString := `

github.com/R9295/schloss v0.1`
	_, err := ParseLockfile(&lockfileString)
	if assert.Error(t, err) {
		assert.Equal(t, err, fmt.Errorf("Invalid entry on line 3"))
	}

}

func TestParseLockfileBadSecondLine(t *testing.T) {
	lockfileString := `
github.com/R9295/schloss v0.1 h1:hash
github.com/R9295/schloss v1.0`
	_, err := ParseLockfile(&lockfileString)
	if assert.Error(t, err) {
		assert.Equal(t, err, fmt.Errorf("Invalid entry on line 3"))
	}

}

func TestParseLockfileNoSecondLine(t *testing.T) {
	lockfileString := `
github.com/R9295/schloss v0.1 h1:hash

github.com/R9295/schloss v1.0`
	_, err := ParseLockfile(&lockfileString)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			err,
			fmt.Errorf("Expected go.mod entry for github.com/R9295/schloss on line 3"),
		)
	}
}

func TestParseLockfileLineNameNotMatching(t *testing.T) {
	lockfileString := `
github.com/R9295/schloss v0.1 h1:hash
github.com/R9295/NotSchloss v1.0 h1:hash`
	_, err := ParseLockfile(&lockfileString)
	if assert.Error(t, err) {
		assert.Equal(
			t,
			err,
			fmt.Errorf(
				"Invalid go.mod entry on line 3. Expected github.com/R9295/schloss got github.com/R9295/NotSchloss",
			),
		)
	}
}
