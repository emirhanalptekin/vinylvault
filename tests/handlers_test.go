package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/emirhanalptekin/vinylvault/internal/api"
	"github.com/emirhanalptekin/vinylvault/internal/db"
	"github.com/emirhanalptekin/vinylvault/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	// Set up router with the health check route
	router := gin.Default()
	router.GET("/", api.HealthCheck)

	// Create a request to send to the above route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	// Assert we got an OK response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

// TestGetAlbums tests the GET /albums endpoint
func TestGetAlbums(t *testing.T) {
	// Set up mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	db.SetDBPool(mock)

	// Define expected rows returned from the database
	rows := mock.NewRows([]string{"id", "title", "artist_id", "artist_name", "release_year", "genre_id", "genre_name", "genre_icon", "notes", "rating", "condition"}).
		AddRow("alb-001", "The Dark Side of the Moon", "art-001", "Pink Floyd", "1973", "gen-001", "Rock", "🎸", "Original pressing", 5, "Excellent")

	// Set up expected query with regular expression for flexibility
	queryRegex := regexp.QuoteMeta(`
		SELECT a.id, a.title, a.artist_id, ar.name, a.release_year, a.genre_id, g.name, g.icon, a.notes, a.rating, a.condition 
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.id
		JOIN genres g ON a.genre_id = g.id
	`)
	mock.ExpectQuery(queryRegex).WillReturnRows(rows)

	// Set up router with the albums route
	router := gin.Default()
	router.GET("/albums", api.GetAlbums)

	// Create a request to send to the above route
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums", nil)
	router.ServeHTTP(w, req)

	// Assert we got an OK response with the expected data
	assert.Equal(t, http.StatusOK, w.Code)

	// Check that the response contains the expected album data
	var albums []models.Album
	err = json.Unmarshal(w.Body.Bytes(), &albums)
	assert.NoError(t, err)
	assert.Len(t, albums, 1)
	assert.Equal(t, "alb-001", albums[0].ID)
	assert.Equal(t, "The Dark Side of the Moon", albums[0].Title)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestGetAlbumByID tests the GET /albums/:id endpoint
func TestGetAlbumByID(t *testing.T) {
	// Set up mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	db.SetDBPool(mock)

	// Set up expected query and rows
	queryRegex := regexp.QuoteMeta(`
		SELECT a.id, a.title, a.artist_id, ar.name, a.release_year, a.genre_id, g.name, g.icon, a.notes, a.rating, a.condition 
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.id
		JOIN genres g ON a.genre_id = g.id
		WHERE a.id = $1
	`)
	rows := mock.NewRows([]string{"id", "title", "artist_id", "artist_name", "release_year", "genre_id", "genre_name", "genre_icon", "notes", "rating", "condition"}).
		AddRow("alb-001", "The Dark Side of the Moon", "art-001", "Pink Floyd", "1973", "gen-001", "Rock", "🎸", "Original pressing", 5, "Excellent")

	mock.ExpectQuery(queryRegex).WithArgs("alb-001").WillReturnRows(rows)

	// Set up router with the album by ID route
	router := gin.Default()
	router.GET("/albums/:id", api.GetAlbumByID)

	// Create a request to send
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/albums/alb-001", nil)
	router.ServeHTTP(w, req)

	// Assert we got an OK response with the expected data
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var album models.Album
	err = json.Unmarshal(w.Body.Bytes(), &album)
	assert.NoError(t, err)
	assert.Equal(t, "alb-001", album.ID)
	assert.Equal(t, "The Dark Side of the Moon", album.Title)
	assert.Equal(t, "Pink Floyd", album.Artist.Name)
	assert.Equal(t, "Rock", album.Genre.Name)

	// Check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestCreateAlbum tests the POST /albums endpoint
func TestCreateAlbum(t *testing.T) {
	// Set up mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	db.SetDBPool(mock)

	// Define the album to create
	album := models.Album{
		ID:          "alb-test",
		Title:       "Test Album",
		ArtistID:    "art-001",
		ReleaseYear: "2023",
		GenreID:     "gen-001",
		Notes:       "Test notes",
		Rating:      4,
		Condition:   models.ConditionMint,
	}

	// Set up expected query
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO albums (id, title, artist_id, release_year, genre_id, notes, rating, condition)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)).WithArgs(
		album.ID, album.Title, album.ArtistID, album.ReleaseYear, album.GenreID, album.Notes, album.Rating, album.Condition,
	).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	// Set up router
	router := gin.Default()
	router.POST("/albums", api.CreateAlbum)

	// Create request body
	jsonValue, _ := json.Marshal(album)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/albums", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert successful creation
	assert.Equal(t, http.StatusCreated, w.Code)

	// Check response
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "alb-test", response["id"])

	// Check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestUpdateAlbum tests the PUT /albums/:id endpoint
func TestUpdateAlbum(t *testing.T) {
	// Set up mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	db.SetDBPool(mock)

	// Define the album to update
	album := models.Album{
		Title:       "Updated Album",
		ArtistID:    "art-001",
		ReleaseYear: "2023",
		GenreID:     "gen-001",
		Notes:       "Updated notes",
		Rating:      5,
		Condition:   models.ConditionExcellent,
	}

	// Set up expected query
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE albums 
		SET title = $2, artist_id = $3, release_year = $4, genre_id = $5, notes = $6, rating = $7, condition = $8
		WHERE id = $1
	`)).WithArgs(
		"alb-001", album.Title, album.ArtistID, album.ReleaseYear, album.GenreID, album.Notes, album.Rating, album.Condition,
	).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Set up router
	router := gin.Default()
	router.PUT("/albums/:id", api.UpdateAlbum)

	// Create request
	jsonValue, _ := json.Marshal(album)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/albums/alb-001", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assert successful update
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Album updated successfully", response["message"])

	// Check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestDeleteAlbum tests the DELETE /albums/:id endpoint
func TestDeleteAlbum(t *testing.T) {
	// Set up mock database
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("Unable to create mock database connection: %v", err)
	}
	defer mock.Close()
	db.SetDBPool(mock)

	// Set up expected query
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM albums WHERE id = $1")).
		WithArgs("alb-001").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Set up router
	router := gin.Default()
	router.DELETE("/albums/:id", api.DeleteAlbum)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/albums/alb-001", nil)
	router.ServeHTTP(w, req)

	// Assert successful deletion
	assert.Equal(t, http.StatusOK, w.Code)

	// Check response
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Album deleted successfully", response["message"])

	// Check expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
