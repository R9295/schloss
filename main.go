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
	diff, _ := diffparser.Parse(GetGitDiff())
	lockFileRegex := regexp.MustCompile(fmt.Sprintf("%s", lockFileStruct.fileName))
	for _, file := range diff.Files {
		// We don't care about deleted (0) or created (2) modes so we only want modified (1)
		if lockFileRegex.MatchString(file.NewName) && file.Mode == 1 {
			fmt.Printf("Lockfile %s edited\n", file.NewName)
		}
	}
}
