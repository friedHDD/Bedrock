package handler

import (
	"github.com/friedHDD/Bedrock/core/library"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"sync"
)

/**
Programming is a theory.......
先锁门来看是不是在cleaning，知道在清理后马上解锁门让下一个人来看
知道没有人在清理之后挂牌"正在清理"然后立即解锁，这样别人就能看到"正在清理"
这样就防止别人和我同时查看，然后我开始清理，他以为没有在清理，然后就冲突了
Gemini:用互斥锁（Mutex）来保护一个共享状态（isCleaning 标志位）的读和写，确保这些操作是“原子性”的。锁的持有时间要尽可能地短，一旦对状态的操作完成，就立刻释放锁，绝不能让耗时的任务长时间占有锁。
*/

var (
	cleanMutex sync.Mutex
	isCleaning = false
)

func LibraryCleanHandler(c *gin.Context) {
	libraryYamlFile := "./data/index/library.yaml"
	//lock and check if the other task running
	cleanMutex.Lock() //锁门
	if isCleaning {   //进去康康
		cleanMutex.Unlock() //有人在清理，解锁门马上逃跑(?
		c.JSON(http.StatusConflict, gin.H{"message": "A clean job is already in progress."})
		return
	}

	//mark as started
	isCleaning = true   //没有人在清理，那我挂牌说明我在清理
	cleanMutex.Unlock() //解锁门，这样别人也能看

	go func() {
		defer func() {
			cleanMutex.Lock()
			isCleaning = false
			cleanMutex.Unlock()
			log.Println("Library clean process finished.")
		}()

		log.Println("Starting library clean process...")

		//read
		yamlFile, err := os.ReadFile(libraryYamlFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Println("Library file not found, nothing to clean.")
				return
			}
			log.Printf("Cleanup failed: could not read %s: %v", libraryYamlFile, err)
			return
		}

		var libraryData library.Data
		if err := yaml.Unmarshal(yamlFile, &libraryData); err != nil {
			log.Printf("Cleanup failed: could not unmarshal %s: %v", libraryYamlFile, err)
			return
		}

		if libraryData.Books == nil || len(libraryData.Books) == 0 {
			log.Println("Library is empty, cleanup not needed.")
			return
		}

		keptBooks := make(map[string]library.BookInfo)
		deletedCount := 0

		for id, book := range libraryData.Books {
			if _, statErr := os.Stat(book.OriginPath); os.IsNotExist(statErr) {
				//this book has passed away
				deletedCount++
				log.Printf("File not found for book '%s' at: %s. Removing record.", book.BookName, book.OriginPath)
			} else if statErr != nil {
				//maybe it is still alive :P
				log.Printf("Could not verify file for book '%s' due to an error: %v. Keeping record.", book.BookName, statErr)
				keptBooks[id] = book
			} else {
				//still alive
				keptBooks[id] = book
			}
		}

		libraryData.Books = keptBooks
		updatedYaml, marshalErr := yaml.Marshal(&libraryData)
		if marshalErr != nil {
			log.Printf("Cleanup failed: could not marshal cleaned library data: %v", marshalErr)
			return
		}

		if writeErr := os.WriteFile(libraryYamlFile, updatedYaml, 0644); writeErr != nil {
			log.Printf("Cleanup failed: could not write updated library file %s: %v", libraryYamlFile, writeErr)
			return
		}
		log.Printf("Cleanup complete. Removed %d non-existent book records.", deletedCount)
	}()

	c.JSON(http.StatusAccepted, gin.H{"message": "Clean job has been accepted."})
}
