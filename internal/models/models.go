package models

// Album represents a vinyl record in the collection
// @Description Information about a vinyl record
type Album struct {
	ID          string         `json:"id" example:"alb-12345678" format:"uuid"`
	Title       string         `json:"title" example:"The Dark Side of the Moon" binding:"required"`
	ArtistID    string         `json:"artist_id" example:"art-001"`
	Artist      *Artist        `json:"artist,omitempty"`
	ReleaseYear string         `json:"release_year" example:"1973" binding:"required"`
	GenreID     string         `json:"genre_id" example:"gen-001"`
	Genre       *Genre         `json:"genre,omitempty"`
	Notes       string         `json:"notes" example:"Original pressing with posters and stickers"`
	Rating      int            `json:"rating" example:"5" minimum:"1" maximum:"5"` // 1-5 stars
	Condition   AlbumCondition `json:"condition" example:"Excellent" enums:"Mint,Excellent,Very Good,Good,Fair,Poor"`
}

// Artist represents a musical artist
// @Description Information about a music artist
type Artist struct {
	ID   string `json:"id" example:"art-001" format:"uuid"`
	Name string `json:"name" example:"Pink Floyd" binding:"required"`
}

// Genre represents a music genre
// @Description Information about a music genre
type Genre struct {
	ID   string `json:"id" example:"gen-001" format:"uuid"`
	Name string `json:"name" example:"Rock" binding:"required"`
	Icon string `json:"icon" example:"🎸"` // Could be an emoji
}

// AlbumCondition represents the physical condition of a vinyl record
// @Description Physical condition of a vinyl record
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
// @Description Standard error response format
type ErrorResponse struct {
	Error string `json:"error" example:"Failed to retrieve album"`
}
