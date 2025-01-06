package main

import (
	"log"
	"overlay-image/controllers"

	"github.com/gin-gonic/gin"
)

// Main initializes and starts the server.
func main() {
	r := gin.Default()

	r.POST("/template", controllers.TemplateController)

	err := r.Run(":8000"); 
	
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Println("Server is runnig on port 8080")
}
