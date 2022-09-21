package gomod

import (
	"fmt"
	"strings"

	"github.com/R9295/schloss/core"
)

type LockfilePackage struct {
	Name            string
	Version         string
	Checksum        string
	ModfileChecksum string
	Dependencies    []string
}

type Lockfile struct {
	Packages []LockfilePackage
}

func destructureLine(line []string) (string, string, string) {
	return line[0], line[1], line[2]
}

func parseLine(line string, lineNum int) (string, string, string, error) {
	lineSplit := strings.Split(line, " ")
	if len(lineSplit) != 3 {
		return "","","", fmt.Errorf("Invalid entry on line %d", lineNum)
	}
	name, version, checksum := destructureLine(lineSplit)
	return name, version, checksum, nil
}

func ParseLockfile(lockfile *string) (Lockfile, error) {
	parsedLockfile := Lockfile{}
	lines := strings.Split(*lockfile, "\n")
	for i := 0; i < len(lines); {
		if lines[i] != "" {
			name, version, checksum, err := parseLine(lines[i], i + 1)
			if err != nil {
				return parsedLockfile, err
			}
			if lines[i+1] == "" {
				return parsedLockfile, fmt.Errorf(
				"Expected go.mod entry for %s on line %d",
				name,
				i+2,
			)
			}
			secondLineName, secondLineVersion, secondLineChecksum, err := parseLine(lines[i+1], i+2)
			if err != nil {
				return parsedLockfile, err
			}
			if name != secondLineName {
				return parsedLockfile, fmt.Errorf(
					"Invalid go.mod entry on line %d. Expected %s got %s",
					i+2,
					name,
					secondLineName,
				)
			} else if secondLineVersion != fmt.Sprintf("%s/go.mod", version) {

				return parsedLockfile, fmt.Errorf(
					"Invalid go.mod entry on line %d. Expected %s got %s",
					i+2,
					name,
					secondLineName,
				)
			}
			parsedLockfile.Packages = append(parsedLockfile.Packages, LockfilePackage{
				Name:            name,
				Version:         version[1:],
				Checksum:        checksum,
				ModfileChecksum: secondLineChecksum,
			})
			i += 2
		} else {
			i++
		}
	}
	return parsedLockfile, nil
}

func Diff(rootFile *string, oldLockfile *string, newLockfile *string, diffList *[]core.Diff) error {
	return nil
}
