package main

import "log"

type LockFileType struct {
	fileName string
	format string
}

func GetLockFileType(lockFileType string) LockFileType {
	switch lockFileType {
	case "yarn":
		return LockFileType{fileName: "yarn.lock", format: "json"}
	case "poetry":
		return LockFileType{fileName: "poetry.lock", format: "toml"}
	case "npm":
		return LockFileType{fileName: "package-lock.json", format: "json"}
	}
	log.Fatalf("Unsupported lockfile type %s", lockFileType)
	return LockFileType{}
}
