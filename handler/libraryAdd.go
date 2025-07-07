package handler

import (
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// BookInfo md5 is the key
type BookInfo struct {
	OriginPath string `yaml:"originPath"`
	Series     string `yaml:"series"`
	BookName   string `yaml:"bookName"`
}

type LibraryData struct {
	Books map[string]BookInfo `yaml:"books"`
}

func LibraryAddHandler(c *gin.Context) {
	bookSeries := "ungrouped"
	libraryResPath := "./data/res/library"
	libraryResUngroupedPath := path.Join(libraryResPath, bookSeries)
	libraryYamlFile := "./data/index/library.yaml"

	//init library ungrouped path
	_err := os.MkdirAll(libraryResUngroupedPath, 0755)
	if _err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": _err.Error()})
		return
	}

	queryFilePath := c.Query("file")
	originPath, err := utils.ConvertPath(queryFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	bookName := filepath.Base(originPath)

	//copy ebook
	newPath := path.Join(libraryResUngroupedPath, bookName)
	err = utils.CopyFile(originPath, newPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	bookMd5Id := utils.Md5(bookSeries + "/" + bookName)

	/**start library init**/
	yamlFile, err := os.ReadFile(libraryYamlFile)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to read %s: %v", libraryYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to read library data"})
		return
	}

	var libraryData LibraryData
	//read data from library
	if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
		log.Printf("Failed to unmarshal %s: %v", libraryYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to parse library data"})
		return
	}

	//if empty
	if libraryData.Books == nil {
		libraryData.Books = make(map[string]BookInfo)
	}
	/**end library init**/

	//if this book existed
	if _, exists := libraryData.Books[bookMd5Id]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "book already existed"})
		return
	}

	//add
	libraryData.Books[bookMd5Id] = BookInfo{
		OriginPath: originPath,
		Series:     bookSeries,
		BookName:   bookName,
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
	if err := os.WriteFile(libraryYamlFile, newYaml, 0644); err != nil {
		log.Printf("Failed to write to %s: %v", libraryYamlFile, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save new data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book added successfully"})
}
