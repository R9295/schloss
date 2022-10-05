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

func TestParseSubDependenciesGraph(t *testing.T) {
	graph := `github.com/R9295/schloss github.com/R9295/parserkiosk@v0.0.1
github.com/R9295/schloss gopkg.in/yaml@v2.2.2
github.com/R9295/parserkiosk@v0.0.1 gopkg.in/ruebezahl@v2.1.3`
	parsed, err := parseSubDependencyGraph(graph)
	assert.Equal(t, err, nil)
	expected := make(map[string][]string, 0)
	expected["github.com/R9295/schloss"] = []string{"github.com/R9295/parserkiosk", "gopkg.in/yaml"}
	expected["github.com/R9295/parserkiosk"] = []string{"gopkg.in/ruebezahl"}
	assert.Equal(t, expected, parsed)
}

func TestParseSubDependencies(t *testing.T) {
	lockfile := Lockfile{
		Packages: []LockfilePackage{
			{
				Name:            "github.com/R9295/schloss",
				Version:         "0.1",
				Checksum:        "h1:hash",
				ModfileChecksum: "h1:hashModFile",
			},
			{
				Name:            "github.com/R9295/parserkiosk",
				Version:         "0.0.1",
				Checksum:        "h1:hash",
				ModfileChecksum: "h1:hashModFile",
			},
			{
				Name:            "gopkg.in/yaml",
				Version:         "0.0.3",
				Checksum:        "h1:hash",
				ModfileChecksum: "h1:hashModFile",
			},
			{
				Name:            "gopkg.in/ruebezahl",
				Version:         "0.0.4",
				Checksum:        "h1:hash",
				ModfileChecksum: "h1:hashModFile",
			},
		},
	}
	graphString := `github.com/R9295/schloss github.com/R9295/parserkiosk@v0.0.1
github.com/R9295/parserkiosk@v0.0.1 gopkg.in/ruebezahl@v0.0.4
github.com/R9295/parserkiosk@v0.0.1 gopkg.in/yaml@v0.0.3`
	parseSubDependencies(&lockfile, graphString)
	assert.Equal(t, lockfile.Packages[0].Dependencies, []string{"github.com/R9295/parserkiosk"})
	assert.Equal(
		t,
		lockfile.Packages[1].Dependencies,
		[]string{"gopkg.in/ruebezahl", "gopkg.in/yaml"},
	)
	assert.Equal(t, lockfile.Packages[2].Dependencies, []string(nil))
	assert.Equal(t, lockfile.Packages[3].Dependencies, []string(nil))
}
