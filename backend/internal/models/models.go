package models

type User struct {
	ID          int64   `json:"id"`
	DisplayName string  `json:"display_name"`
	Handle      string  `json:"handle"`
	CNZBalance  float64 `json:"cnz_balance"`
	CreatedAt   string  `json:"created_at"`
}

type Artwork struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	Caption   string    `json:"caption"`
	CreatedAt string    `json:"created_at"`
	Artist    *User     `json:"artist,omitempty"`
	Comments  []Comment `json:"comments,omitempty"`
}

type Comment struct {
	ID        int64   `json:"id"`
	ArtworkID int64   `json:"artwork_id"`
	UserID    int64   `json:"user_id"`
	RatingCNZ int     `json:"rating_cnz"`
	Body      string  `json:"body"`
	CreatedAt string  `json:"created_at"`
	Author    *User   `json:"author,omitempty"`
	ArtistCNZ float64 `json:"artist_cnz"`
}

type LedgerEntry struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	ArtworkID *int64  `json:"artwork_id,omitempty"`
	Amount    float64 `json:"amount"`
	Direction string  `json:"direction"`
	Reason    string  `json:"reason"`
	CreatedAt string  `json:"created_at"`
}
