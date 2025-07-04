package functions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileListItem struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	LastModify string `json:"lastModify"`
	Permission string `json:"permission"`
}

func ListDirectoryHandler(c *gin.Context) {
	// folder=
	folderPath := c.Query("folder")
	if folderPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder query parameter is required"})
		return
	}

	if strings.HasPrefix(folderPath, "~") {
		homeDir, err := os.UserHomeDir() // ~ means homedir
		if err != nil {
			log.Printf("Error getting user home directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not resolve home directory"})
			return
		}
		folderPath = filepath.Join(homeDir, folderPath[1:]) // generate true path
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", folderPath, err)

		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("folder not found: %s", folderPath)})
		} else if os.IsPermission(err) {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("permission denied: %s", folderPath)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read directory"})
		}
		return
	}

	var fileList []FileListItem
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			// if file disappeared
			log.Printf("Could not get info for file %s: %v", entry.Name(), err)
			continue
		}

		itemType := "file"
		if info.IsDir() {
			itemType = "folder"
		}

		fileList = append(fileList, FileListItem{
			Name:       info.Name(),
			Type:       itemType,
			LastModify: info.ModTime().Format(time.RFC3339), // use ISO8601
			Permission: fmt.Sprintf("%04o", info.Mode().Perm()),
		})
	}

	c.JSON(http.StatusOK, fileList)
}
