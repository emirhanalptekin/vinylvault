package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emirhanalptekin/vinylvault/internal/api"
	"github.com/emirhanalptekin/vinylvault/internal/config"
	"github.com/emirhanalptekin/vinylvault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupRouter initializes the router with all routes for testing
func setupRouter() *gin.Engine {
	// Use a test database URL - this could be pointing to a test container
	_ = &config.Config{
		DatabaseUrl: "postgres://vinylvault:vinylvault@localhost:5432/vinylvault_test?sslmode=disable",
		Port:        "8080",
	}

	// Skip actual database connection during tests
	// In a real scenario, you might use a test database or container
	// db.InitializeDB(testConfig.DatabaseUrl)

	// Set up the router with all routes
	router := gin.Default()
	api.RegisterRoutes(router)

	return router
}

// TestIntegrationHealthCheck tests the health check endpoint integration
func TestIntegrationHealthCheck(t *testing.T) {
	// Skip if not running integration tests
	t.Skip("Integration test - skipped in normal unit test runs")

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

// TestIntegrationAlbumCRUD tests the complete CRUD flow for albums
func TestIntegrationAlbumCRUD(t *testing.T) {
	// Skip if not running integration tests
	t.Skip("Integration test - skipped in normal unit test runs")

	router := setupRouter()

	// Step 1: Create a new album
	album := models.Album{
		Title:       "Integration Test Album",
		ArtistID:    "art-001", // assuming this exists in test DB
		ReleaseYear: "2023",
		GenreID:     "gen-001", // assuming this exists in test DB
		Notes:       "Created during integration test",
		Rating:      5,
		Condition:   models.ConditionMint,
	}

	jsonValue, _ := json.Marshal(album)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/albums", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &createResponse)
	assert.NoError(t, err)
	albumID := createResponse["id"]
	assert.NotEmpty(t, albumID)

	// Step 2: Get the created album
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/albums/"+albumID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedAlbum models.Album
	err = json.Unmarshal(w.Body.Bytes(), &retrievedAlbum)
	assert.NoError(t, err)
	assert.Equal(t, album.Title, retrievedAlbum.Title)

	// Step 3: Update the album
	retrievedAlbum.Title = "Updated Integration Test Album"
	retrievedAlbum.Notes = "Updated during integration test"

	jsonValue, _ = json.Marshal(retrievedAlbum)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/albums/"+albumID, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 4: Get the updated album
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/albums/"+albumID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedAlbum models.Album
	err = json.Unmarshal(w.Body.Bytes(), &updatedAlbum)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Integration Test Album", updatedAlbum.Title)
	assert.Equal(t, "Updated during integration test", updatedAlbum.Notes)

	// Step 5: Delete the album
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/albums/"+albumID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Step 6: Verify the album was deleted
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/albums/"+albumID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestIntegrationGetArtists tests retrieving all artists
func TestIntegrationGetArtists(t *testing.T) {
	// Skip if not running integration tests
	t.Skip("Integration test - skipped in normal unit test runs")

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/artists", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var artists []models.Artist
	err := json.Unmarshal(w.Body.Bytes(), &artists)
	assert.NoError(t, err)
	assert.NotEmpty(t, artists)
}

// TestIntegrationGetGenres tests retrieving all genres
func TestIntegrationGetGenres(t *testing.T) {
	// Skip if not running integration tests
	t.Skip("Integration test - skipped in normal unit test runs")

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/genres", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var genres []models.Genre
	err := json.Unmarshal(w.Body.Bytes(), &genres)
	assert.NoError(t, err)
	assert.NotEmpty(t, genres)
}
