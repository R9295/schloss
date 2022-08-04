package cargo

import (
	"github.com/R9295/schloss/contrib/toml"
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

func ParseLockfiles(oldLockfile string, newLockfile string) (Lockfile, Lockfile) {
	var newLockfileToml Lockfile
	var oldLockfileToml Lockfile
	toml.DecodeToml(newLockfile, &newLockfileToml)
	toml.DecodeToml(oldLockfile, &oldLockfileToml)
	return oldLockfileToml, newLockfileToml
}


func DiffLockfiles(oldLockfileToml *Lockfile, newLockfileToml *Lockfile) []core.Diff {
	diffList := make([]core.Diff, 0)
	return diffList
}
