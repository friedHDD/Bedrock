package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitAll() error {
	const basePath = "./data"
	folderList := []string{"index", "res/library", "res/music"}
	for _, folder := range folderList {
		if err := os.MkdirAll(filepath.Join(basePath, folder), 0755); err != nil {
			return fmt.Errorf("failed to create %s: %v", folder, err)
		}
	}
	return nil
}
