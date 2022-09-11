package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/R9295/schloss/contrib/cargo"
	"github.com/R9295/schloss/contrib/poetry"
	"github.com/R9295/schloss/core"
	"github.com/jessevdk/go-flags"
	"github.com/waigani/diffparser"
)

var args struct {
	LockfileType    string `short:"t" long:"type" description:"Type of lockfile" required:"true"`
	IgnoreUntracked bool   `long:"ignore-untracked" description:"Ignore Untracked Log Files"`
	Format          string `short:"f" long:"fmt" default:"human" description:"Format of output, options: json, human."`
	CommitAmount    uint   `long:"commit-amount" default:"1" description:"diff commit amount (HEAD~commitAmount)."`
	Log             bool   `long:"log" description:"Log your lockfile diff"`
	LogFile         string `long:"log-file" default:"schloss.log" description:"File to log your diff into."`
	OverrideLog     bool   `long:"override-log" description:"Override latest log entry if it's a duplicate"`
}

func run() error {
	_, err := flags.Parse(&args)
	if err != nil {
		return nil
	}
	lockfileType, err := core.GetLockfileType(args.LockfileType)
	if err != nil {
		return err
	}
	fmt.Println("----------------------------------------------------------------------")
	fmt.Printf("Running schloss for type: %s\n", args.LockfileType)
	fmt.Println("----------------------------------------------------------------------")
	start := time.Now()
	if !args.IgnoreUntracked {
		untrackedLockfiles, amount, err := core.CheckUntrackedFiles(lockfileType.FileName)
		if err != nil {
			return err
		}
		if amount > 0 {
			fmt.Println("Error: You have untracked lockfiles. Please add them to source control.")
			for _, file := range untrackedLockfiles {
				fmt.Println(file)
			}
			fmt.Println(
				"If you think this is a bug, you can silence it with --ignore-untracked and file a bug report",
			)
			return nil
		}
	}
	gitDiff, err := core.GetAllDiff(args.CommitAmount)
	if err != nil {
		return err
	}
	parsedDiff, err := diffparser.Parse(gitDiff)
	if err != nil {
		return err
	}
	lockfileRegex := regexp.MustCompile(lockfileType.FileName)
	for _, file := range parsedDiff.Files {
		if lockfileRegex.MatchString(file.NewName) && file.Mode == diffparser.MODIFIED {
			fmt.Printf("Lockfile %s modified\n", file.NewName)
			lockfileDiff, err := core.GetSingleDiff(file.NewName, args.CommitAmount)
			if err != nil {
				return err
			}
			parsedLockfileDiff, err := diffparser.Parse(lockfileDiff)
			if err != nil {
				return err
			}
			diffedFile := parsedLockfileDiff.Files[0]
			newLockfile, oldLockfile := core.GetLockfilesFromDiff(diffedFile)
			rootFile, err := core.GetRootFile(file.NewName, lockfileType.RootFile)
			if err != nil {
				return err
			}
			var diffList []core.Diff
			if args.LockfileType == "poetry" {
				if err := poetry.Diff(&rootFile, &oldLockfile, &newLockfile, &diffList); err != nil {
					return err
				}
			} else {
				if err := cargo.Diff(&rootFile, &oldLockfile, &newLockfile, &diffList); err != nil {
					return err
				}
			}
			var rendered string
			if args.Format == "json" {
				rendered, err = core.RenderJSON(&diffList)
				if err != nil {
					return err
				}
			} else {
				rendered = core.RenderHumanReadable(&diffList)
			}
			if args.Log {
				if err := core.Log(args.LogFile, rendered, args.OverrideLog); err != nil {
					return err
				}
			}
			fmt.Println(rendered)
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
