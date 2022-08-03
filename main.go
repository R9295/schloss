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
		Name string
		Version string
		Dependencies map[string]interface{}
	}
	Metadata struct {
		PythonVersions string `toml:"python-versions"`
		ContentHash string `toml:"content-hash"`
		LockVersion string `toml:"lock-version"`
		Files map[string][]PoetryLockfileMetadataFile
	}
}


func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Printf("Running schloss for %s type: %s\n", opts.LockfilePath, opts.LockfileType)
	fmt.Println("----------------------------------------------------------------------")
	lockFileStruct := GetLockFileType(opts.LockfileType)
	if !opts.IgnoreUntracked {
		untrackedLockfiles, amount := CheckUntrackedFiles(lockFileStruct.fileName)
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
	lockFileRegex := regexp.MustCompile(lockFileStruct.fileName)
	for _, file := range gitDiff.Files {
		if lockFileRegex.MatchString(file.NewName) && file.Mode == diffparser.MODIFIED {
			fmt.Printf("Lockfile %s modified\n", file.NewName)
			lockfileDiff, _ := diffparser.Parse(GetSingleDiff(file.NewName))
			newFile := ""
			oldFile := ""
			file = lockfileDiff.Files[0]
			for _, hunk := range file.Hunks {
				for _, line := range hunk.WholeRange.Lines {
					switch line.Mode {
					case diffparser.ADDED:
						newFile = fmt.Sprintf("%s\n%s", newFile, line.Content)
					case diffparser.REMOVED:
						oldFile = fmt.Sprintf("%s\n%s", oldFile, line.Content)
					case diffparser.UNCHANGED:
						oldFile = fmt.Sprintf("%s\n%s", oldFile, line.Content)
						newFile = fmt.Sprintf("%s\n%s", newFile, line.Content)
					}
				}
			}
			var newFileToml PoetryLockfile
			var oldFileToml PoetryLockfile
			decodeToml(newFile, &newFileToml)
			decodeToml(oldFile, &oldFileToml)
			changelog, err := diff.Diff(oldFileToml, newFileToml)
			for _, item := range changelog{
				fmt.Println(item.Path)
				fmt.Println(item.Type)
				fmt.Println()
			}
		}
	}
}
