package cargo

import (
	"github.com/R9295/schloss/core"
)

type LockfilePackage struct {
	Name         string
	Version      string
	Source       string
	Checksum     string
	Dependencies []string
}

type Lockfile struct {
	Package []LockfilePackage
}

func DiffLockfiles(oldLockfileToml *Lockfile, newLockfileToml *Lockfile) []core.Diff {
	diffList := make([]core.Diff, 0)
	return diffList
}
