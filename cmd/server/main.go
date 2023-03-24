package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Set up your API routes
	// ...

	// Start the server
	r.Run(":8080")
}
