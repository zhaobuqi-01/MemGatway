package main

import (
	_ "gateway/init" // Import the init package to initialize configuration and logging
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// Set up your API routes
	// ...

	// Start the server
	r.Run(":8080")
}
