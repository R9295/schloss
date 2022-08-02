package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	LockfileType string `short:"t" long:"type" description:"Type of lockfile" required:"true"`
	LockfilePath string `short:"p" long:"path" description:"Path to lockfile" required:"true"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}
	fmt.Printf("Running schloss for %s type: %s\n", opts.LockfilePath, opts.LockfileType)
}
