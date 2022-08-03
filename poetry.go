package main

import (
	"fmt"

	"github.com/r3labs/diff/v3"
)

type PoetryLockfileMetadataFile struct {
	File string
	Hash string
}
type PoetryLockfile struct {
	Package []struct {
		Name         string
		Version      string
		Dependencies map[string]interface{}
	}
	Metadata struct {
		PythonVersions string `toml:"python-versions"`
		ContentHash    string `toml:"content-hash"`
		LockVersion    string `toml:"lock-version"`
		Files          map[string][]PoetryLockfileMetadataFile
	}
}

func poetryGetPackageVersion(lockfileDiff *diff.Changelog, index string) string {
	version := ""
	for _, change := range *lockfileDiff {
		if change.Path[0] == "Package" && change.Path[1] == index && change.Path[2] == "Version" {
			version = fmt.Sprintf("%s", change.To)
			break
		}
	}
	return version
}

func poetryCheckNewPackages(lockfileDiff *diff.Changelog) {
	for _, change := range *lockfileDiff {
		if change.Type == diff.CREATE {
			if change.Path[0] == "Package" && change.Path[2] == "Name" {
				version := poetryGetPackageVersion(lockfileDiff, change.Path[1])
				fmt.Println(fmt.Sprintf("add %s at %s", change.To, version))
			}
		}
	}
}
