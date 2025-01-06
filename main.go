package main

import (
	"log"
	"overlay-image/controllers"

	"github.com/gin-gonic/gin"
	// "github.com/aws/aws-lambda-go/events"
  // "github.com/aws/aws-lambda-go/lambda"
)

// Main initializes and starts the server.
func main() {
	r := gin.Default()

	r.POST("/template", controllers.TemplateController)

	err := r.Run(":8080"); 
	
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	log.Println("Server is runnig on port 8080")
}
