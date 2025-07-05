package handler

import (
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFileHandler(c *gin.Context) {

	queryFilePath := c.Query("file")
	filePath, err := utils.ConvertPath(queryFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not access file"})
		}
		return
	}

	if info.IsDir() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "damn it is a directory"})
		return
	}

	fileName := filepath.Base(filePath)
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Header("Content-Type", "application/octet-stream")

	c.File(filePath)
}
