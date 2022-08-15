package core

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

func CheckUntrackedFiles(fileName string) ([]string, int) {
	untrackedFiles, err := exec.Command("git", "ls-files", "--others", "--exclude-standard").Output()
	if err != nil {
		log.Fatal(err.Error())
	}
	reUntracked := regexp.MustCompile(fmt.Sprintf(`.*%s`, fileName))
	untrackedLockfiles := reUntracked.FindAllString(string(untrackedFiles), -1)
	return untrackedLockfiles, len(untrackedLockfiles)
}

func GetAllDiff(commitAmount uint) string {
	diff, err := exec.Command("git", "diff", fmt.Sprintf("HEAD~%d", commitAmount)).Output()
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(diff)
}

func GetSingleDiff(filePath string, commitAmount uint) string {
	fmt.Println(filePath)

	// your diff line length better not be bigger than that number! TODO: handle if not
	diff, err := exec.Command("git", "diff", fmt.Sprintf("HEAD~%d", commitAmount), "-U99999999999999999", filePath).Output()
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(diff)
}
