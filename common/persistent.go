package common

import (
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
