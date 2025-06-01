package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	fmt.Println("🚀 Server starting on http://localhost:8080 ...")

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		fmt.Println("❌ Error starting server:", err)
	}
}
