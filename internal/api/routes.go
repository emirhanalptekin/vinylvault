package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up the API routes
func RegisterRoutes(router *gin.Engine) {
	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
	}))

	// Health check
	router.GET("/", HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Albums routes
	router.GET("/albums", GetAlbums)
	router.GET("/albums/:id", GetAlbumByID)
	router.POST("/albums", CreateAlbum)
	router.PUT("/albums/:id", UpdateAlbum)
	router.DELETE("/albums/:id", DeleteAlbum)

	// Artists route
	router.GET("/artists", GetArtists)

	// Genres route
	router.GET("/genres", GetGenres)
}
