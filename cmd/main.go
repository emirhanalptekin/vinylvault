package main

import (
	"log"

	"github.com/emirhanalptekin/vinylvault/internal/api"
	"github.com/emirhanalptekin/vinylvault/internal/config"
	"github.com/emirhanalptekin/vinylvault/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.GetAppConfig("internal/config/config.yml")

	// Initialize database connection
	db.InitializeDB(cfg.DatabaseUrl)

	// Set up Gin router
	router := gin.Default()

	// Register API routes
	api.RegisterRoutes(router)

	// Start the server
	log.Printf("Starting server on port %s...\n", cfg.Port)
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
