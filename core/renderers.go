package core

import (
	"encoding/json"
	"fmt"
)

func RenderHumanReadable(diffList *[]Diff) {
	for _, item := range *diffList {
		fmt.Println(item.RenderHumanReadable())
	}
}

func RenderJSON(diffList *[]Diff) error {
	jsonDiff, err := json.Marshal(diffList)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonDiff))
	return nil
}
