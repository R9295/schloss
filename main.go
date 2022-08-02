package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/jessevdk/go-flags"
	"github.com/waigani/diffparser"
)

var opts struct {
	LockfileType    string `short:"t" long:"type" description:"Type of lockfile" required:"true"`
	LockfilePath    string `short:"p" long:"path" description:"Path to lockfile" required:"true"`
	IgnoreUntracked bool   `long:"ignore-untracked" description:"Ignore Untracked Log Files"`
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
		untrackedLogFiles, amount := CheckUntrackedFiles(lockFileStruct.fileName)
		if amount > 0 {
			fmt.Println("Error: You have untracked lockfiles. Please add them to source control.")
			for _, file := range untrackedLogFiles {
				fmt.Println(file)
			}
			fmt.Println("If you think this is a bug, you can silence it with --ignore-untracked and file a bug report")
			return
		}
	}
	diff, _ := diffparser.Parse(GetAllDiff())
	lockFileRegex := regexp.MustCompile(lockFileStruct.fileName)
	for _, file := range diff.Files {
		// We don't care about deleted (0) or created (2) modes so we only want modified (1)
		if lockFileRegex.MatchString(file.NewName) && file.Mode == 1 {
			fmt.Printf("Lockfile %s edited\n", file.NewName)
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
			fmt.Println(newFile)
		}
	}
}
