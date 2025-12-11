package main

import (
	"log"
	"os"

	"simple-pdf-converter/handlers"
	"simple-pdf-converter/middleware"
	"simple-pdf-converter/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Initialize PDFium
	if err := utils.InitPDFium(); err != nil {
		log.Fatalf("Failed to initialize PDFium: %v", err)
	}
	defer utils.ClosePDFium()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.Default()

	// API routes with middleware
	api := router.Group("/api")
	api.Use(middleware.APIKeyAuth())
	{
		api.POST("/convert", handlers.ConvertPDF)
	}

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
