package handler

import (
	"github.com/friedHDD/Bedrock/core/library"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

func LibraryGetHandler(c *gin.Context) {
	bookID := c.Param("id")
	if bookID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Book ID is required"})
		return
	}

	libraryYamlFile := "./data/index/library.yaml"

	yamlFile, err := os.ReadFile(libraryYamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Library index not found"})
			return
		}
		log.Printf("Failed to read %s: %v", libraryYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read library data"})
		return
	}

	var libraryData library.Data
	if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
		log.Printf("Failed to unmarshal %s: %v", libraryYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse library data"})
		return
	}

	//id->book
	bookInfo, exists := libraryData.Books[bookID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	//check if file exists
	if _, statErr := os.Stat(bookInfo.OriginPath); os.IsNotExist(statErr) {
		log.Printf("File not found for book '%s' at: %s.", bookInfo.BookName, bookInfo.OriginPath)
		c.JSON(http.StatusNotFound, gin.H{"message": "Book file not found on server"})
		return
	}

	c.File(bookInfo.OriginPath)
}
