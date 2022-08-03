package main

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
