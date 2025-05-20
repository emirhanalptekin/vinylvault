package api

import (
	"net/http"

	"github.com/emirhanalptekin/vinylvault/internal/db"
	"github.com/emirhanalptekin/vinylvault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Health check endpoint
// @Summary Health check
// @Description Check if the API is running
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetAlbums handles GET /albums request
// @Summary Get all albums
// @Description Retrieve all albums in the collection
// @Tags albums
// @Produce json
// @Success 200 {array} models.Album
// @Failure 500 {object} models.ErrorResponse
// @Router /albums [get]
func GetAlbums(c *gin.Context) {
	albums, err := db.GetAlbums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve albums"})
		return
	}
	c.JSON(http.StatusOK, albums)
}

// GetAlbumByID handles GET /albums/:id request
// @Summary Get album by ID
// @Description Retrieve a specific album by its ID
// @Tags albums
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {object} models.Album
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /albums/{id} [get]
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
// @Summary Create a new album
// @Description Add a new album to the collection
// @Tags albums
// @Accept json
// @Produce json
// @Param album body models.Album true "Album Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /albums [post]
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
// @Summary Update an album
// @Description Update an existing album's information
// @Tags albums
// @Accept json
// @Produce json
// @Param id path string true "Album ID"
// @Param album body models.Album true "Album Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /albums/{id} [put]
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
// @Summary Delete an album
// @Description Remove an album from the collection
// @Tags albums
// @Produce json
// @Param id path string true "Album ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} models.ErrorResponse
// @Router /albums/{id} [delete]
func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")

	if err := db.DeleteAlbum(id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to delete album"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// GetArtists handles GET /artists request
// @Summary Get all artists
// @Description Retrieve all artists in the collection
// @Tags artists
// @Produce json
// @Success 200 {array} models.Artist
// @Failure 500 {object} models.ErrorResponse
// @Router /artists [get]
func GetArtists(c *gin.Context) {
	artists, err := db.GetArtists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve artists"})
		return
	}
	c.JSON(http.StatusOK, artists)
}

// GetGenres handles GET /genres request
// @Summary Get all genres
// @Description Retrieve all music genres in the collection
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 500 {object} models.ErrorResponse
// @Router /genres [get]
func GetGenres(c *gin.Context) {
	genres, err := db.GetGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve genres"})
		return
	}
	c.JSON(http.StatusOK, genres)
}
