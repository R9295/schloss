package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/jessevdk/go-flags"
	diff "github.com/r3labs/diff/v3"
	"github.com/waigani/diffparser"
)

var opts struct {
	LockfileType    string `short:"t" long:"type" description:"Type of lockfile" required:"true"`
	LockfilePath    string `short:"p" long:"path" description:"Path to lockfile" required:"true"`
	IgnoreUntracked bool   `long:"ignore-untracked" description:"Ignore Untracked Log Files"`
}

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

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Running schloss for %s type: %s\n", opts.LockfilePath, opts.LockfileType)
	fmt.Println("----------------------------------------------------------------------")
	lockfileStruct := GetLockfileType(opts.LockfileType)
	if !opts.IgnoreUntracked {
		untrackedLockfiles, amount := CheckUntrackedFiles(lockfileStruct.fileName)
		if amount > 0 {
			fmt.Println("Error: You have untracked lockfiles. Please add them to source control.")
			for _, file := range untrackedLockfiles {
				fmt.Println(file)
			}
			fmt.Println("If you think this is a bug, you can silence it with --ignore-untracked and file a bug report")
			return
		}
	}
	gitDiff, _ := diffparser.Parse(GetAllDiff())
	lockfileRegex := regexp.MustCompile(lockfileStruct.fileName)
	for _, file := range gitDiff.Files {
		if lockfileRegex.MatchString(file.NewName) && file.Mode == diffparser.MODIFIED {
			fmt.Printf("Lockfile %s modified\n", file.NewName)
			lockfileDiff, _ := diffparser.Parse(GetSingleDiff(file.NewName))
			newLockfile := ""
			oldLockfile := ""
			file = lockfileDiff.Files[0]
			GetLockfileFromDiff(&newLockfile, &oldLockfile, file)
			var newLockfileToml PoetryLockfile
			var oldLockfileToml PoetryLockfile
			DecodeToml(newLockfile, &newLockfileToml)
			DecodeToml(oldLockfile, &oldLockfileToml)
			changelog, err := diff.Diff(oldLockfileToml, newLockfileToml)
			if err != nil {
				log.Fatal(err)
			}
			for _, item := range changelog {
				fmt.Println(item.Path)
				fmt.Println(item.Type)
				fmt.Println()
			}
		}
	}
}
