package handler

import (
	"fmt"
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

type (
	fileDetails struct {
		TruePath   string `json:"truePath"`
		Size       int64  `json:"size"`
		LastModify string `json:"lastModify"`
		Permission string `json:"permission"`
	}
)

func FileDetailHandler(c *gin.Context) {

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not access file info"})
		}
		return
	}

	fileDetailData := fileDetails{
		TruePath:   filePath,
		Size:       info.Size(),
		LastModify: info.ModTime().Format(time.RFC3339),
		Permission: fmt.Sprintf("%04o", info.Mode().Perm()),
	}

	c.JSON(http.StatusOK, fileDetailData)
}
