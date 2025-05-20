package models

// Album represents a vinyl record in the collection
type Album struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	ArtistID    string         `json:"artist_id"`
	Artist      *Artist        `json:"artist,omitempty"`
	ReleaseYear string         `json:"release_year"`
	GenreID     string         `json:"genre_id"`
	Genre       *Genre         `json:"genre,omitempty"`
	Notes       string         `json:"notes"`
	Rating      int            `json:"rating"` // 1-5 stars
	Condition   AlbumCondition `json:"condition"`
}

// Artist represents a musical artist
type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Genre represents a music genre
type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"` // Could be an emoji
}

// AlbumCondition represents the physical condition of a vinyl record
type AlbumCondition string

const (
	ConditionMint      AlbumCondition = "Mint"
	ConditionExcellent AlbumCondition = "Excellent"
	ConditionVeryGood  AlbumCondition = "Very Good"
	ConditionGood      AlbumCondition = "Good"
	ConditionFair      AlbumCondition = "Fair"
	ConditionPoor      AlbumCondition = "Poor"
)

// ErrorResponse standardizes error responses
type ErrorResponse struct {
	Error string `json:"error"`
}
