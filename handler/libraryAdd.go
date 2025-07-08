package handler

import (
	"github.com/friedHDD/Bedrock/core/library"
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func LibraryAddHandler(c *gin.Context) {
	bookSeries := "ungrouped"
	libraryResPath := "./data/res/library"
	libraryResUngroupedPath := path.Join(libraryResPath, bookSeries)

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

	bookMd5Id := library.Md5(bookSeries + "/" + bookName)

	bookToAdd := make(map[string]library.BookInfo)
	bookToAdd[bookMd5Id] = library.BookInfo{
		OriginPath: originPath,
		Series:     bookSeries,
		BookName:   bookName,
	}

	err, num := library.Add(bookToAdd)
	if err != nil {
		log.Printf("Failed to add book to library index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if num == 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "book already existed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book added successfully"})
}
