package api

import (
	"net/http"

	"github.com/emirhanalptekin/vinylvault/internal/db"
	"github.com/emirhanalptekin/vinylvault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetAlbums handles GET /albums request
func GetAlbums(c *gin.Context) {
	albums, err := db.GetAlbums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve albums"})
		return
	}
	c.JSON(http.StatusOK, albums)
}

// GetAlbumByID handles GET /albums/:id request
func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	album, err := db.GetAlbumByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve album"})
		return
	}

	if album == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Album not found"})
		return
	}

	c.JSON(http.StatusOK, album)
}

// CreateAlbum handles POST /albums request
func CreateAlbum(c *gin.Context) {
	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid album data"})
		return
	}

	// Generate a UUID if not provided
	if album.ID == "" {
		album.ID = "alb-" + uuid.New().String()[:8]
	}

	if err := db.CreateAlbum(album); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create album"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": album.ID})
}

// UpdateAlbum handles PUT /albums/:id request
func UpdateAlbum(c *gin.Context) {
	id := c.Param("id")

	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid album data"})
		return
	}

	// Ensure the ID in the path matches the ID in the body
	album.ID = id

	if err := db.UpdateAlbum(album); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to update album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album updated successfully"})
}

// DeleteAlbum handles DELETE /albums/:id request
func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")

	if err := db.DeleteAlbum(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// GetArtists handles GET /artists request
func GetArtists(c *gin.Context) {
	artists, err := db.GetArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve artists"})
		return
	}
	c.JSON(http.StatusOK, artists)
}

// GetGenres handles GET /genres request
func GetGenres(c *gin.Context) {
	genres, err := db.GetGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve genres"})
		return
	}
	c.JSON(http.StatusOK, genres)
}
