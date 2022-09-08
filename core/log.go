package core

import (
	"bufio"
	"fmt"
	"os"
)

const ENTRY_SEPARATOR string = "----------------------------------------------------------------------------------------------------"

func checkForDuplicateEntry(fileName string, header string) (bool, error) {
	logFile, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer logFile.Close()
	scanner := bufio.NewScanner(logFile)
	scanner.Scan()
	return scanner.Text() == header, scanner.Err()
}

func getLatestEntryLineCount(fileName string) (int, error) {
	lineCount := 0
	logFile, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return lineCount, err
	}
	defer logFile.Close()
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		lineCount++
		if scanner.Text() == ENTRY_SEPARATOR {
			lineCount++
			break
		}
	}
	return lineCount, scanner.Err()
}

func Log(fileName string, diff string, override bool) error {
	if fileName == "" {
		fileName = "schloss.log"
	}
	tempFileName := fmt.Sprintf("%s.temp", fileName)
	commitHash, err := GetPrevCommitHash()
	if err != nil {
		return err
	}
	header := fmt.Sprintf("previous commit: %s", commitHash)
	diff = fmt.Sprintf("%s\n%s\n%s\n", header, diff, ENTRY_SEPARATOR)
	hasDuplicate, err := checkForDuplicateEntry(fileName, header)
	if err != nil {
		return err
	}
	ignoreLines := 0
	if hasDuplicate {
		if override {
			ignoreLines, err = getLatestEntryLineCount(fileName)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("A schloss entry with the previous commit already exists. Use --override-log to override it")
		}
	}
	logFile, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	tempFile, err := os.OpenFile(tempFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()
	defer tempFile.Close()
	if _, err := tempFile.WriteString(diff); err != nil {
		return err
	}
	scanner := bufio.NewScanner(logFile)
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
