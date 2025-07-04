package server

import (
	"flag"
	"fmt"
	"github.com/friedHDD/Bedrock/functions"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start() {
	port := flag.String("port", "9090", "the port for the server to listen on")

	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/list", func(c *gin.Context) {
		functions.ListDirectoryHandler(c)
	})

	address := fmt.Sprintf(":%s", *port)
	log.Printf("Bedrock starting, listening on %s", address)

	if err := r.Run(address); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
