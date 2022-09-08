package toml

import (
	"github.com/BurntSushi/toml"
)

func DecodeToml[T any](text string, lockfileStruct T) error {
	_, err := toml.Decode(text, lockfileStruct)
	if err != nil {
		return err
	}
	return nil
}
func ParseToml[T any](lockfile string) (T, error) {
	var parsedLockfile T
	if err := DecodeToml(lockfile, &parsedLockfile); err != nil {
		return parsedLockfile, err
	}
	return parsedLockfile, nil
}
