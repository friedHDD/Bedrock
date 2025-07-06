package handler

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
)

type bookData struct {
	ID       string `json:"id"`
	Series   string `json:"series"`
	BookName string `json:"bookName"`
}

func LibraryListHandler(c *gin.Context) {
	libraryYamlFile := "./data/index/library.yaml"

	//read
	yamlFile, err := os.ReadFile(libraryYamlFile)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusOK, gin.H{"books": []bookData{}})
			return
		}
		log.Printf("Failed to read %s: %v", libraryYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read library data"})
		return
	}

	var libraryData LibraryData
	if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
		log.Printf("Failed to unmarshal %s: %v", libraryYamlFile, err)
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
			Series:   book.Series,
			BookName: book.BookName,
		})
	}

	c.JSON(http.StatusOK, gin.H{"books": booksResponse})
}
