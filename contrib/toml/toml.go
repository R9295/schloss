package toml

import (
	"log"

	"github.com/BurntSushi/toml"
)

func DecodeToml[T any](text string, lockfileStruct T) {
	_, err := toml.Decode(text, lockfileStruct)
	if err != nil {
		log.Fatal(err)
	}
}
func ParseLockfiles[T any](oldLockfile string, newLockfile string) (T, T) {
	var newLockfileToml T
	var oldLockfileToml T
	DecodeToml(newLockfile, &newLockfileToml)
	DecodeToml(oldLockfile, &oldLockfileToml)
	return oldLockfileToml, newLockfileToml
}
