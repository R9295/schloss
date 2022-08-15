package core

import (
	"fmt"
	"os/exec"
	"regexp"
)

func CheckUntrackedFiles(fileName string) ([]string, int) {
	untrackedFiles, err := exec.Command("git", "ls-files", "--others", "--exclude-standard").Output()
	if err != nil {
		panic(err)
	}
	reUntracked := regexp.MustCompile(fmt.Sprintf(`.*%s`, fileName))
	untrackedLockfiles := reUntracked.FindAllString(string(untrackedFiles), -1)
	return untrackedLockfiles, len(untrackedLockfiles)
}

func GetAllDiff() string {
	diff, err := exec.Command("git", "diff", "HEAD").Output()
	if err != nil {
		panic(err)
	}
	return string(diff)
}

func GetSingleDiff(filePath string) string {
	fmt.Println(filePath)
	// better hope your diff is not larger than the number! TODO: handle incase it is
	diff, err := exec.Command("git", "diff", "HEAD", "-U999999999999999999999", filePath).Output()
	if err != nil {
		panic(err)
	}
	return string(diff)
}
