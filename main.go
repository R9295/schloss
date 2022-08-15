package main

import (
	"encoding/json"
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
	IgnoreUntracked bool   `long:"ignore-untracked" description:"Ignore Untracked Log Files"`
	Format          string `short:"f" long:"fmt" description:"Format of output, options: json, human. Default: human"`
	CommitAmount    uint `long:"commit-amount" description:"diff commit amount (HEAD~commitAmount). Default: 1"`
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
	fmt.Printf("Running schloss for type: %s\n", opts.LockfileType)
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
	var commitAmount uint
	if opts.CommitAmount > 0 {
		commitAmount = opts.CommitAmount
	} else {
		commitAmount = 1
	}
	gitDiff, _ := diffparser.Parse(core.GetAllDiff(commitAmount))
	lockfileRegex := regexp.MustCompile(lockfileType.FileName)
	for _, file := range gitDiff.Files {
		if lockfileRegex.MatchString(file.NewName) && file.Mode == diffparser.MODIFIED {
			fmt.Printf("Lockfile %s modified\n", file.NewName)
			lockfileDiff, _ := diffparser.Parse(core.GetSingleDiff(file.NewName, commitAmount))
			newLockfile := ""
			oldLockfile := ""
			file = lockfileDiff.Files[0]
			core.GetLockfileFromDiff(&newLockfile, &oldLockfile, file)
			var diffList []core.Diff
			if opts.LockfileType == "poetry" {
				oldLockfileToml, err := toml.ParseLockfile[poetry.Lockfile](oldLockfile)
				if err != nil {
					return err
				}
				newLockfileToml, err := toml.ParseLockfile[poetry.Lockfile](newLockfile)
				if err != nil {
					return err
				}
				diffList = poetry.DiffLockfiles(&oldLockfileToml, &newLockfileToml)
			} else {
				oldLockfileToml, err := toml.ParseLockfile[cargo.Lockfile](oldLockfile)
				if err != nil {
					return err
				}
				newLockfileToml, err := toml.ParseLockfile[cargo.Lockfile](newLockfile)
				if err != nil {
					return err
				}
				diffList = cargo.DiffLockfiles(&oldLockfileToml, &newLockfileToml)
			}
			if opts.Format == "json" {
				jsonDiff, err := json.Marshal(diffList)
				if err != nil {
					return err
				}
				fmt.Println(string(jsonDiff))
			} else {
				for _, item := range diffList {
					fmt.Println(item.RenderHumanReadable())
				}
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
