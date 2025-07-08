package handler

import (
	"fmt"
	"github.com/friedHDD/Bedrock/core/library"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	scanMutex  sync.Mutex
	isScanning bool
)

func LibraryScanHandler(c *gin.Context) {
	libraryResPath := "./data/res/library"
	//lock and check if the other task running
	scanMutex.Lock()
	if isScanning {
		scanMutex.Unlock()
		c.JSON(http.StatusConflict, gin.H{"message": "A scan job is already in progress."})
		return
	}

	//mark as started
	isScanning = true
	scanMutex.Unlock()

	go func() {
		defer func() {
			scanMutex.Lock()
			isScanning = false
			scanMutex.Unlock()
			log.Println("Library scan process finished.")
		}()

		log.Println("Starting library scan process...")

		bookToAdd := make(map[string]library.BookInfo)

		seriesDirs, err := os.ReadDir(libraryResPath)
		if err != nil {
			log.Printf("Error reading library directory %s: %v", libraryResPath, err)
			return
		}

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

					bookMd5Id := library.Md5(relativePath)

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

		err, num := library.Add(bookToAdd)
		if err != nil {
			log.Printf("Failed to add book to library index: %v", err)
			return
		}

		message := fmt.Sprintf("Scan complete: Added %d new books.", num)
		log.Println(message)
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Scan job has been accepted."})
}
