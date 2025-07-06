package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type bookData struct {
	ID       string `json:"id"`
	FilePath string `json:"filePath"`
	BookName string `json:"bookName"`
}

func LibraryListHandler(c *gin.Context) {
	const libraryPath = "./data/library"
	booksYamlFile := filepath.Join(libraryPath, "books.yaml")

	//read
	yamlFile, err := os.ReadFile(booksYamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"books": []bookData{}})
			return
		}
		log.Printf("Failed to read %s: %v", booksYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read library data"})
		return
	}

	var libraryData LibraryData
	if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
		log.Printf("Failed to unmarshal %s: %v", booksYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to parse library data"})
		return
	}

	//if books.yaml is empty
	if libraryData.Books == nil {
		c.JSON(http.StatusOK, gin.H{"books": []bookData{}})
		return
	}

	//get
	booksResponse := make([]bookData, 0, len(libraryData.Books))
	for id, book := range libraryData.Books {
		booksResponse = append(booksResponse, bookData{
			ID:       id,
			FilePath: book.FilePath,
			BookName: book.BookName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"books": booksResponse})
}
