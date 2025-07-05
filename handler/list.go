package handler

import (
	"fmt"
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type (
	FileListItem struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		LastModify string `json:"lastModify"`
		Permission string `json:"permission"`
	}
)

func ListDirectoryHandler(c *gin.Context) {

	queryFolderPath := c.Query("folder")
	folderPath, err := utils.ConvertPath(queryFolderPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
