package handler

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// BookInfo md5 is the key
type BookInfo struct {
	FilePath string `yaml:"filePath"`
	BookName string `yaml:"bookName"`
}

type LibraryData struct {
	Books map[string]BookInfo `yaml:"books"`
}

func LibraryAddHandler(c *gin.Context) {
	const libraryPath = "./data/library"
	booksYamlFile := filepath.Join(libraryPath, "books.yaml")

	err := utils.LibraryInit()
	if err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "library cannot be initialized"})
		return
	}

	queryFilePath := c.Query("file")
	filePath, err := utils.ConvertPath(queryFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	filePathMd5 := md5Hash(filePath)
	bookName := filepath.Base(filePath)

	yamlFile, err := os.ReadFile(booksYamlFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to read %s: %v", booksYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read library data"})
		return
	}

	var libraryData LibraryData
	//read data from library
	if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
		log.Printf("Failed to unmarshal %s: %v", booksYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to parse library data"})
		return
	}

	//if empty
	if libraryData.Books == nil {
		libraryData.Books = make(map[string]BookInfo)
	}

	//if this book existed
	if _, exists := libraryData.Books[filePathMd5]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "book already existed"})
		return
	}

	//add
	libraryData.Books[filePathMd5] = BookInfo{
		FilePath: filePath,
		BookName: bookName,
	}

	//pack new data
	var newYaml []byte
	newYaml, err = yaml.Marshal(&libraryData)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to pack data before saving"})
		return
	}

	//write
	if err := os.WriteFile(booksYamlFile, newYaml, 0644); err != nil {
		log.Printf("Failed to write to %s: %v", booksYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save new data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book added successfully"})
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
