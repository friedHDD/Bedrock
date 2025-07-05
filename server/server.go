package server

import (
	"github.com/friedHDD/Bedrock/handler"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start() {

	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/list", func(c *gin.Context) {
		handler.ListDirectoryHandler(c)
	})

	r.GET("/api/download", func(c *gin.Context) {
		handler.DownloadFileHandler(c)
	})

	log.Printf("Bedrock starting, listening on :9090")

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
