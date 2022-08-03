package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/waigani/diffparser"
)

type LockFileType struct {
	fileName string
	format   string
}

func GetLockfileType(lockfileType string) LockFileType {
	switch lockfileType {
	case "yarn":
		return LockFileType{fileName: "yarn.lock", format: "json"}
	case "poetry":
		return LockFileType{fileName: "poetry.lock", format: "toml"}
	case "npm":
		return LockFileType{fileName: "package-lock.json", format: "json"}
	}
	log.Fatalf("Unsupported lockfile type %s", lockfileType)
	return LockFileType{}
}

func DecodeToml(text string, lockfileStruct *PoetryLockfile) {
	_, err := toml.Decode(text, lockfileStruct)
	if err != nil {
		log.Fatal(err)
	}
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
