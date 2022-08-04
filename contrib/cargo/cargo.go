package cargo

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
