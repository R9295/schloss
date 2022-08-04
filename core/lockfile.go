package core

import (
	"fmt"
	"log"

	"github.com/waigani/diffparser"
)

type LockFileType struct {
	FileName string
	Format   string
}

func GetLockfileType(lockfileType string) LockFileType {
	switch lockfileType {
	case "yarn":
		return LockFileType{FileName: "yarn.lock", Format: "json"}
	case "poetry":
		return LockFileType{FileName: "poetry.lock", Format: "toml"}
	case "npm":
		return LockFileType{FileName: "package-lock.json", Format: "json"}
	}
	log.Fatalf("Unsupported lockfile type %s", lockfileType)
	return LockFileType{}
}

func GetLockfileFromDiff(newLockfile *string, oldLockfile *string, lockfileDiffFile *diffparser.DiffFile) {
	for _, hunk := range lockfileDiffFile.Hunks {
		for _, line := range hunk.WholeRange.Lines {
			switch line.Mode {
			case diffparser.ADDED:
				*newLockfile = fmt.Sprintf("%s\n%s", *newLockfile, line.Content)
			case diffparser.REMOVED:
				*oldLockfile = fmt.Sprintf("%s\n%s", *oldLockfile, line.Content)
			case diffparser.UNCHANGED:
				*oldLockfile = fmt.Sprintf("%s\n%s", *oldLockfile, line.Content)
				*newLockfile = fmt.Sprintf("%s\n%s", *newLockfile, line.Content)
			}
		}
	}
}
