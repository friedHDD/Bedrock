package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func LibraryInit() error {
	const libraryPath = "./data/library"
	booksYamlFile := filepath.Join(libraryPath, "books.yaml")
	if err := os.MkdirAll(libraryPath, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %v", libraryPath, err)
	}

	if _, err := os.Stat(booksYamlFile); os.IsNotExist(err) {
		initialContent := []byte("# This file stores books data.\n")
		if err := os.WriteFile(booksYamlFile, initialContent, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %v", booksYamlFile, err)
		}

	} else if err != nil {
		return fmt.Errorf("%s: %v", booksYamlFile, err)
	} else {
		return nil
		//log.Printf("Books data file %s already exists.", booksYamlFile)
	}
	return nil
}
