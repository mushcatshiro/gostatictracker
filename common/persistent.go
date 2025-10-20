package common

import (
	"errors"
	"fmt"
	"os"
)

func PersistToFileSystem(targetDirectory, fname, fileContent string) error {
	file, err := os.Create(targetDirectory + "/" + fname) // truncates if exists
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", fname, err)
	}
	defer file.Close()
	_, err = file.Write([]byte(fileContent))
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}

func ValidateInputDirectory(targetDirectory string) error {
	info, err := os.Stat(targetDirectory)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("Directory specified does not exists: %s", targetDirectory)
	}
	if err != nil {
		return fmt.Errorf("Error during directory check %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Expects a valid path instead got %s", targetDirectory)
	}
	return nil
}

func ValidateInputFilePath(targetFilePath string) error {
	info, err := os.Stat(targetFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("File path specified does not exists: %s", targetFilePath)
	} else if err != nil {
		fmt.Printf("Error accessing file '%s': %v\n", targetFilePath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("Expects a valid file path instead got %s", targetFilePath)
	}
	return nil
}
