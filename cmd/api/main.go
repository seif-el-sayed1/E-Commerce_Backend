package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in .env
	}

	server.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Hello World",
		})
	})

	server.Run(":" + port)
}
