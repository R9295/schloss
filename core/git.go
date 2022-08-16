package core

import (
	"fmt"
	"os/exec"
	"regexp"
)

func CheckUntrackedFiles(fileName string) ([]string, int, error) {
	untrackedFiles, err := exec.Command("git", "ls-files", "--others", "--exclude-standard").Output()
	if err != nil {
		return []string{}, 0, err
	}
	reUntracked := regexp.MustCompile(fmt.Sprintf(`.*%s`, fileName))
	untrackedLockfiles := reUntracked.FindAllString(string(untrackedFiles), -1)
	return untrackedLockfiles, len(untrackedLockfiles), nil
}

func GetAllDiff(commitAmount uint) (string, error) {
	diff, err := exec.Command("git", "diff", fmt.Sprintf("HEAD~%d", commitAmount)).Output()
	if err != nil {
		return "", err
	}
	return string(diff), nil
}

func GetSingleDiff(filePath string, commitAmount uint) (string, error) {
	fmt.Println(filePath)

	// your diff line length better not be bigger than that number! TODO: handle if not
	diff, err := exec.Command("git", "diff", fmt.Sprintf("HEAD~%d", commitAmount), "-U99999999999999999", filePath).Output()
	if err != nil {
		return "", err
	}
	return string(diff), nil
}
