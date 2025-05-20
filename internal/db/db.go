package db

import (
	"context"
	"log"
	"time"

	"github.com/emirhanalptekin/vinylvault/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool is an interface for database operations, useful for mocking in tests
type DBPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Close()
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// Database connection pool
var dbPool DBPool

// InitializeDB initializes the database connection pool
func InitializeDB(connString string) {
	var err error
	pool, err := pgxpool.New(context.Background(), connString)
	for attempts := 0; attempts < 5; attempts++ {
		if err == nil {
			SetDBPool(pool)
			log.Println("Connected to the database successfully.")
			return
		}
		log.Printf("Unable to connect to database, retrying in 3 seconds... (%d/5)\n", attempts+1)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("Unable to connect to database: %v\n", err)
}

// Function to set the dbPool variable - mostly used for testing
func SetDBPool(pool DBPool) {
	dbPool = pool
}

// GetAlbums retrieves all albums from the database
func GetAlbums() ([]models.Album, error) {
	rows, err := dbPool.Query(context.Background(), `
		SELECT a.id, a.title, a.artist_id, ar.name, a.release_year, a.genre_id, g.name, g.icon, a.notes, a.rating, a.condition 
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.id
		JOIN genres g ON a.genre_id = g.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []models.Album
	for rows.Next() {
		var album models.Album
		var artistName, genreName, genreIcon string

		err = rows.Scan(
			&album.ID,
			&album.Title,
			&album.ArtistID,
			&artistName,
			&album.ReleaseYear,
			&album.GenreID,
			&genreName,
			&genreIcon,
			&album.Notes,
			&album.Rating,
			&album.Condition,
		)
		if err != nil {
			return nil, err
		}

		album.Artist = &models.Artist{ID: album.ArtistID, Name: artistName}
		album.Genre = &models.Genre{ID: album.GenreID, Name: genreName, Icon: genreIcon}

		albums = append(albums, album)
	}

	return albums, nil
}

// GetAlbumByID retrieves a single album by ID
func GetAlbumByID(id string) (*models.Album, error) {
	var album models.Album
	var artistName, genreName, genreIcon string

	err := dbPool.QueryRow(context.Background(), `
		SELECT a.id, a.title, a.artist_id, ar.name, a.release_year, a.genre_id, g.name, g.icon, a.notes, a.rating, a.condition 
		FROM albums a
		JOIN artists ar ON a.artist_id = ar.id
		JOIN genres g ON a.genre_id = g.id
		WHERE a.id = $1
	`, id).Scan(
		&album.ID,
		&album.Title,
		&album.ArtistID,
		&artistName,
		&album.ReleaseYear,
		&album.GenreID,
		&genreName,
		&genreIcon,
		&album.Notes,
		&album.Rating,
		&album.Condition,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No album found
		}
		return nil, err
	}

	album.Artist = &models.Artist{ID: album.ArtistID, Name: artistName}
	album.Genre = &models.Genre{ID: album.GenreID, Name: genreName, Icon: genreIcon}

	return &album, nil
}

// CreateAlbum adds a new album to the database
func CreateAlbum(album models.Album) error {
	_, err := dbPool.Exec(context.Background(), `
		INSERT INTO albums (id, title, artist_id, release_year, genre_id, notes, rating, condition)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, album.ID, album.Title, album.ArtistID, album.ReleaseYear, album.GenreID, album.Notes, album.Rating, album.Condition)

	return err
}

// UpdateAlbum updates an existing album
func UpdateAlbum(album models.Album) error {
	_, err := dbPool.Exec(context.Background(), `
		UPDATE albums 
		SET title = $2, artist_id = $3, release_year = $4, genre_id = $5, notes = $6, rating = $7, condition = $8
		WHERE id = $1
	`, album.ID, album.Title, album.ArtistID, album.ReleaseYear, album.GenreID, album.Notes, album.Rating, album.Condition)

	return err
}

// DeleteAlbum removes an album from the database
func DeleteAlbum(id string) error {
	_, err := dbPool.Exec(context.Background(), "DELETE FROM albums WHERE id = $1", id)
	return err
}

// GetArtists retrieves all artists from the database
func GetArtists() ([]models.Artist, error) {
	rows, err := dbPool.Query(context.Background(), "SELECT id, name FROM artists")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artists []models.Artist
	for rows.Next() {
		var artist models.Artist
		err = rows.Scan(&artist.ID, &artist.Name)
		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, nil
}

// GetGenres retrieves all genres from the database
func GetGenres() ([]models.Genre, error) {
	rows, err := dbPool.Query(context.Background(), "SELECT id, name, icon FROM genres")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		err = rows.Scan(&genre.ID, &genre.Name, &genre.Icon)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}

	return genres, nil
}
