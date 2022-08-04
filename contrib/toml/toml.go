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
