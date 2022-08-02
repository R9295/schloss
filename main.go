package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
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
	if !opts.IgnoreUntracked {
		untrackedLogFiles, amount := CheckUntrackedFiles("test.test")
		if amount > 0 {
			fmt.Println("Error: You have untracked lockfiles. Please add them to source control.")
			for _, file := range untrackedLogFiles {
				fmt.Println(file)
			}
			fmt.Println("If you think this is a bug, you can silence it with --ignore-untracked and file a bug report")
			return
		}
	}
}
