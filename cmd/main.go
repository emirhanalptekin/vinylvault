package main

import (
	"log"

	"github.com/emirhanalptekin/vinylvault/internal/api"
	"github.com/emirhanalptekin/vinylvault/internal/config"
	"github.com/emirhanalptekin/vinylvault/internal/db"
	"github.com/gin-gonic/gin"

	_ "github.com/emirhanalptekin/vinylvault/docs"
)

// @title VinylVault API
// @version 1.0
// @description A REST API for managing vinyl record collections
// @host localhost:8080
// @BasePath /
// @schemes http
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
	log.Printf("Swagger UI available at http://localhost:%s/swagger/index.html\n", cfg.Port)
	err := router.Run(":" + cfg.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
