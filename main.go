package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/R9295/schloss/contrib/cargo"
	"github.com/R9295/schloss/contrib/poetry"
	"github.com/R9295/schloss/contrib/toml"
	"github.com/R9295/schloss/core"
	"github.com/jessevdk/go-flags"
	"github.com/waigani/diffparser"
)

var opts struct {
	LockfileType    string `short:"t" long:"type" description:"Type of lockfile" required:"true"`
	LockfilePath    string `short:"p" long:"path" description:"Path to lockfile" required:"true"`
	IgnoreUntracked bool   `long:"ignore-untracked" description:"Ignore Untracked Log Files"`
}

func run() error {
	_, err := flags.Parse(&opts)
	if err != nil {
		return err
	}
	lockfileType, err := core.GetLockfileType(opts.LockfileType)
	if err != nil {
		return err
	}
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("Running schloss for %s type: %s\n", opts.LockfilePath, opts.LockfileType)
	fmt.Println("----------------------------------------------------------------------")
	start := time.Now()
	if !opts.IgnoreUntracked {
		untrackedLockfiles, amount := core.CheckUntrackedFiles(lockfileType.FileName)
		if amount > 0 {
			fmt.Println("Error: You have untracked lockfiles. Please add them to source control.")
			for _, file := range untrackedLockfiles {
				fmt.Println(file)
			}
			fmt.Println("If you think this is a bug, you can silence it with --ignore-untracked and file a bug report")
			return nil
		}
	}
	gitDiff, _ := diffparser.Parse(core.GetAllDiff())
	lockfileRegex := regexp.MustCompile(lockfileType.FileName)
	for _, file := range gitDiff.Files {
		if lockfileRegex.MatchString(file.NewName) && file.Mode == diffparser.MODIFIED {
			fmt.Printf("Lockfile %s modified\n", file.NewName)
			lockfileDiff, _ := diffparser.Parse(core.GetSingleDiff(file.NewName))
			newLockfile := ""
			oldLockfile := ""
			file = lockfileDiff.Files[0]
			core.GetLockfileFromDiff(&newLockfile, &oldLockfile, file)
			var diffList []core.Diff
			if opts.LockfileType == "poetry" {
				oldLockfileToml, newLockfileToml := toml.ParseLockfiles[poetry.Lockfile](oldLockfile, newLockfile)
				diffList = poetry.DiffLockfiles(&oldLockfileToml, &newLockfileToml)
			} else {
				oldLockfileToml, newLockfileToml := toml.ParseLockfiles[cargo.Lockfile](oldLockfile, newLockfile)
				diffList = cargo.DiffLockfiles(&oldLockfileToml, &newLockfileToml)
			}
			for _, item := range diffList {
				fmt.Println(item.RenderHumanReadable())
			}
		}
	}
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("Time elapsed: %s", time.Since(start)))
	fmt.Println("----------------------------------------------------------------------")
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
