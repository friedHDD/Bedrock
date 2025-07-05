package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ConvertPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("damn query parameter is required")
	}
	if strings.HasPrefix(path, "/~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("can not get user home: %v", err)
		}
		return filepath.Join(homeDir, path[2:]), nil
	}
	fmt.Printf(path)

	return path, nil
}
