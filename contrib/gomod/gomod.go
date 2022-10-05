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
		return "", "", "", fmt.Errorf("Invalid entry on line %d", lineNum)
	}
	name, version, checksum := destructureLine(lineSplit)
	return name, version, checksum, nil
}

func ParseLockfile(lockfile *string) (Lockfile, error) {
	parsedLockfile := Lockfile{}
	lines := strings.Split(*lockfile, "\n")
	for lineNum := 0; lineNum < len(lines); {
		if lines[lineNum] != "" {
			name, version, checksum, err := parseLine(lines[lineNum], lineNum+1)
			if err != nil {
				return parsedLockfile, err
			}
			if lines[lineNum+1] == "" {
				return parsedLockfile, fmt.Errorf(
					"Expected go.mod entry for %s on line %d",
					name,
					lineNum+2,
				)
			}
			secondLineName, secondLineVersion, secondLineChecksum, err := parseLine(lines[lineNum+1], lineNum+2)
			if err != nil {
				return parsedLockfile, err
			}
			if name != secondLineName {
				return parsedLockfile, fmt.Errorf(
					"Invalid go.mod entry on line %d. Expected %s got %s",
					lineNum+2,
					name,
					secondLineName,
				)
			} else if secondLineVersion != fmt.Sprintf("%s/go.mod", version) {
				return parsedLockfile, fmt.Errorf(
					"Invalid go.mod entry on line %d. Expected %s got %s",
					lineNum+2,
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
			lineNum += 2
		} else {
			lineNum++
		}
	}
	return parsedLockfile, nil
}

func removeVersionFromDepString(dep string) string {
	return strings.Split(dep, "@")[0]
}

func parseSubDependencyGraph(subDeps string) (map[string][]string, error) {
	graph := make(map[string][]string, 0)
	lines := strings.Split(subDeps, "\n")
	for lineNum, line := range lines {
		splitLine := strings.Split(line, " ")
		if len(splitLine) != 2 {
			return graph, fmt.Errorf(
				"Invalid \"go mod graph\" line entry. Line number %d",
				lineNum+1,
			)
		}
		dep := removeVersionFromDepString(splitLine[0])
		subDep := removeVersionFromDepString(splitLine[1])
		graph[dep] = append(graph[dep], subDep)

	}
	return graph, nil
}

func parseSubDependencies(lockfile *Lockfile, graph string) error {
	subDeps, err := parseSubDependencyGraph(string(graph))
	if err != nil {
		return err
	}
	for index, dep := range lockfile.Packages {
		for _, subDep := range subDeps[dep.Name] {
			lockfile.Packages[index].Dependencies = append(
				lockfile.Packages[index].Dependencies,
				subDep,
			)
		}
	}
	return nil
}
func Diff(rootFile *string, oldLockfile *string, newLockfile *string, diffList *[]core.Diff) error {
	return nil
}
