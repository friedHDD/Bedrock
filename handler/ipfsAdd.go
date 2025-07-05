package handler

import (
	"fmt"
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	"net/http"
	"os"
)

func IPFSAddHandler(c *gin.Context) {

	sh := shell.NewShell("localhost:5001")
	queryFilePath := c.Query("file")
	filePath, err := utils.ConvertPath(queryFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "file cannot be read"})
		}
		return
	}
	defer file.Close()

	cidFromFile, err := sh.Add(file)
	if err != nil {
		errMessage := fmt.Sprintf("error adding  file: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cid": cidFromFile})
}
