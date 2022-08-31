package core

import (
	"bufio"
	"fmt"
	"os"
)

func Log(fileName string, diff string) error {
	if fileName == "" {
		fileName = "schloss.log"
	}
	tempFileName := fmt.Sprintf("%s.temp", fileName)
	logFile, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	tempFile, err := os.OpenFile(tempFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()
	defer tempFile.Close()
	commitHash, err := GetPrevCommitHash()
	if err != nil {
		return err
	}
	diff = fmt.Sprintf("previous commit: %s\n%s\n\n", commitHash, diff)
	if _, err := tempFile.WriteString(diff); err != nil {
		return err
	}
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		if _, err := tempFile.WriteString(fmt.Sprintf("%s\n", scanner.Text())); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	err = os.Rename(tempFileName, fileName)
	if err != nil {
		return err
	}
	return nil
}
