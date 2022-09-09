package core

import (
	"encoding/json"
	"fmt"
)

func RenderHumanReadable(diffList *[]Diff) string {
	var text string
	for index, item := range *diffList {
		text = fmt.Sprintf("%s\n%d. %s", text, index+1, item.RenderHumanReadable())
	}
	return text
}

func RenderJSON(diffList *[]Diff) (string, error) {
	jsonDiff, err := json.Marshal(diffList)
	if err != nil {
		return "", err
	}
	return string(jsonDiff), nil
}
