package handler

import (
	"github.com/friedHDD/Bedrock/core/library"
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LibraryScanHandler(c *gin.Context) {
	libraryResPath := "./data/res/library"

	seriesDirs, err := os.ReadDir(libraryResPath)
	if err != nil {
		log.Printf("Error reading library directory %s: %v", libraryResPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not read library directory"})
		return
	}

	bookToAdd := make(map[string]library.BookInfo)

	//enter folders of each series
	for _, seriesDir := range seriesDirs {
		if !seriesDir.IsDir() || seriesDir.Name() == "ungrouped" {
			continue
		}

		series := seriesDir.Name()
		seriesPath := filepath.Join(libraryResPath, series)

		walkErr := filepath.WalkDir(seriesPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".epub") {
				relativePath, relErr := filepath.Rel(libraryResPath, path)
				if relErr != nil {
					log.Printf("Could not create relative path for %s: %v", path, relErr)
					relativePath = path
				}

				bookMd5Id := utils.Md5(relativePath)

				//a book shouldn't appear in bookToAdd multiple times
				if _, exists := bookToAdd[bookMd5Id]; exists {
					return nil
				}

				bookToAdd[bookMd5Id] = library.BookInfo{
					OriginPath: path,
					Series:     series,
					BookName:   d.Name(),
				}
			}
			return nil
		})

		if walkErr != nil {
			log.Printf("Error scanning series folder %s: %v", seriesPath, walkErr)
			continue
		}

	}

	var num int
	err, num = library.Add(bookToAdd)
	if err != nil {
		log.Printf("Failed to add book to library index: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	var message string
	if num == 0 {
		message = "No new books found"
	} else {
		message = "Added " + strconv.Itoa(num) + " books"
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}
