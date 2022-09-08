package core

import (
	"bufio"
	"fmt"
	"os"
)

const ENTRY_SEPARATOR string = "----------------------------------------------------------------------------------------------------"

func Log(fileName string, diff string, override bool) error {
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
	header := fmt.Sprintf("previous commit: %s", commitHash)
	diff = fmt.Sprintf("%s\n%s\n%s\n", header, diff, ENTRY_SEPARATOR)
	if _, err := tempFile.WriteString(diff); err != nil {
		return err
	}
	scanner := bufio.NewScanner(logFile)
	scanner.Scan()
	ignoreLines := 0
	if scanner.Text() == header {
		if override {
			for scanner.Scan() {
				ignoreLines++
				if scanner.Text() == ENTRY_SEPARATOR {
					ignoreLines++
					break
				}
			}
		} else {
			if err := os.Remove(tempFileName); err != nil {
				return err
			}
			return fmt.Errorf("A schloss entry with the previous commit already exists. Use --override-log to override it")
		}
	}
	logFile.Seek(0, 0)
	for scanner.Scan() {
		if ignoreLines == 0 {
			if _, err := tempFile.WriteString(fmt.Sprintf("%s\n", scanner.Text())); err != nil {
				return err
			}
		} else {
			ignoreLines--
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
