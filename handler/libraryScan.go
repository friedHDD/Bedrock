package handler

import (
	"github.com/friedHDD/Bedrock/utils"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func LibraryScanHandler(c *gin.Context) {
	libraryResPath := "./data/res/library"
	libraryYamlFile := "./data/index/library.yaml"

	seriesDirs, err := os.ReadDir(libraryResPath)
	if err != nil {
		log.Printf("Error reading library directory %s: %v", libraryResPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not read library directory"})
		return
	}
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

	newBooks := 0

	//enter folders of each series
	for _, seriesDir := range seriesDirs {
		if !seriesDir.IsDir() || seriesDir.Name() == "ungrouped" {
			continue
		}

		seriesName := seriesDir.Name()
		seriesPath := filepath.Join(libraryResPath, seriesName)

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

				if _, exists := libraryData.Books[bookMd5Id]; exists {
					return nil
				}

				bookName := d.Name()
				libraryData.Books[bookMd5Id] = BookInfo{
					OriginPath: path,
					Series:     seriesName,
					BookName:   bookName,
				}

				var newYaml []byte
				newYaml, err = yaml.Marshal(&libraryData)
				if err != nil {
					log.Printf("Failed to marshal data: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to pack data before saving"})
					return nil
				}

				//write
				if err := os.WriteFile(libraryYamlFile, newYaml, 0644); err != nil {
					log.Printf("Failed to write to %s: %v", libraryYamlFile, err)
					c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to save new data"})
					return nil
				}

				newBooks++
			}
			return nil
		})

		if walkErr != nil {
			log.Printf("Error scanning series folder %s: %v", seriesPath, walkErr)
			continue
		}

	}
	var message string
	if newBooks == 0 {
		message = "No new books found"
	} else {
		message = "Added " + strconv.Itoa(newBooks) + " books"
	}
	c.JSON(http.StatusOK, gin.H{"message": message})
}
