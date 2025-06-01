package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	fmt.Println("🚀 Server starting on http://localhost:8080 ...")

	router.Run(":8080")
}
