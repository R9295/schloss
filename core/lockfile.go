package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/waigani/diffparser"
)

type LockFileType struct {
	FileName string
	Format   string
	RootFile string
}

func GetLockfileType(lockfileType string) (LockFileType, error) {
	switch lockfileType {
	case "cargo":
		return LockFileType{FileName: "Cargo.lock", Format: "toml", RootFile: "Cargo.toml"}, nil
	case "poetry":
		return LockFileType{
			FileName: "poetry.lock",
			Format:   "toml",
			RootFile: "pyproject.toml",
		}, nil
	case "go":
		return LockFileType{
			FileName: "go.sum",
			Format: "custom",
			RootFile: "go.mod",
		}, nil
	}
	return LockFileType{}, fmt.Errorf("cli: lockfile type \"%s\" is not supported.", lockfileType)
}

func GetLockfilesFromDiff(lockfileDiffFile *diffparser.DiffFile) (string, string) {
	newLockfile := ""
	oldLockfile := ""
	for _, hunk := range lockfileDiffFile.Hunks {
		for _, line := range hunk.WholeRange.Lines {
			switch line.Mode {
			case diffparser.ADDED:
				newLockfile = fmt.Sprintf("%s\n%s", newLockfile, line.Content)
			case diffparser.REMOVED:
				oldLockfile = fmt.Sprintf("%s\n%s", oldLockfile, line.Content)
			case diffparser.UNCHANGED:
				oldLockfile = fmt.Sprintf("%s\n%s", oldLockfile, line.Content)
				newLockfile = fmt.Sprintf("%s\n%s", newLockfile, line.Content)
			}
		}
	}
	return newLockfile, oldLockfile
}

func GetRootFile(lockfilePath string, rootFileName string) (string, error) {
	rootFilePathSplit := strings.Split(lockfilePath, "/")
	rootFilePathSplit = rootFilePathSplit[:len(rootFilePathSplit)-1]
	rootFilePath := fmt.Sprintf("%s/%s", strings.Join(rootFilePathSplit, "/"), rootFileName)
	rootFile, err := os.ReadFile(rootFilePath)
	if err != nil {
		return "", err
	}
	rootFileStr := string(rootFile)
	return rootFileStr, nil
}
