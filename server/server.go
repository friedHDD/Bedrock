package server

import (
	"github.com/friedHDD/Bedrock/handler"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

func Start() {

	r := gin.Default()
	r.StaticFS("/assets", http.Dir("./dist/assets"))

	api := r.Group("/api")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		api.GET("/list", func(c *gin.Context) {
			handler.ListDirectoryHandler(c)
		})

		api.GET("/download", func(c *gin.Context) {
			handler.DownloadHandler(c)
		})

		api.GET("/file/details", func(c *gin.Context) {
			handler.FileDetailHandler(c)
		})

		api.GET("/ipfs/add", func(c *gin.Context) {
			handler.IPFSAddHandler(c)
		})
	}
	r.NoRoute(func(c *gin.Context) {
		indexPath := filepath.Join("./dist", "index.html")
		c.File(indexPath)
	})
	log.Printf("Bedrock starting, listening on :9090")

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
