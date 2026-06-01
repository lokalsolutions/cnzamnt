package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cnzamnt/backend/internal/models"
)

const StartingCNZ = 5000
const ArtistEarnRate = 0.10

var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")
var ErrInsufficientCNZ = errors.New("not enough CNZ")

type Service struct {
	db *sql.DB
}

type CreateUserInput struct {
	DisplayName string `json:"display_name"`
	Handle      string `json:"handle"`
}

type CreateArtworkInput struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Caption  string `json:"caption"`
}

type CreateCommentInput struct {
	RatingCNZ int    `json:"rating_cnz"`
	Body      string `json:"body"`
}

func New(database *sql.DB) *Service {
	return &Service{db: database}
}

func (s *Service) CreateUser(input CreateUserInput) (models.User, error) {
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Handle = strings.TrimSpace(input.Handle)
	if input.DisplayName == "" {
		return models.User{}, fmt.Errorf("display_name is required")
	}
	if input.Handle == "" {
		return models.User{}, fmt.Errorf("handle is required")
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT INTO users (display_name, handle, cnz_balance) VALUES (?, ?, ?)`, input.DisplayName, input.Handle, StartingCNZ)
	if err != nil {
		return models.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.User{}, err
	}
	if _, err = tx.Exec(`INSERT INTO cnz_ledger (user_id, amount, direction, reason) VALUES (?, ?, ?, ?)`,
		id, StartingCNZ, "credit", "signup_bonus"); err != nil {
		return models.User{}, err
	}
	if err = tx.Commit(); err != nil {
		return models.User{}, err
	}
	return s.User(id)
}

func (s *Service) User(id int64) (models.User, error) {
	var user models.User
	err := s.db.QueryRow(`SELECT id, display_name, handle, cnz_balance, created_at FROM users WHERE id = ?`, id).
		Scan(&user.ID, &user.DisplayName, &user.Handle, &user.CNZBalance, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return user, ErrNotFound
	}
	return user, err
}

func (s *Service) CreateArtwork(userID int64, input CreateArtworkInput) (models.Artwork, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.ImageURL = strings.TrimSpace(input.ImageURL)
	input.Caption = strings.TrimSpace(input.Caption)
	if input.Title == "" {
		return models.Artwork{}, fmt.Errorf("title is required")
	}
	if input.ImageURL == "" {
		return models.Artwork{}, fmt.Errorf("image_url is required")
	}
	if _, err := s.User(userID); err != nil {
		return models.Artwork{}, err
	}

	result, err := s.db.Exec(`INSERT INTO artworks (user_id, title, image_url, caption) VALUES (?, ?, ?, ?)`,
		userID, input.Title, input.ImageURL, input.Caption)
	if err != nil {
		return models.Artwork{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return models.Artwork{}, err
	}
	return s.Artwork(id)
}

func (s *Service) Artworks() ([]models.Artwork, error) {
	rows, err := s.db.Query(`SELECT id, user_id, title, image_url, caption, created_at FROM artworks ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, err
	}

	artworks := []models.Artwork{}
	for rows.Next() {
		artwork, err := scanArtwork(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		artworks = append(artworks, artwork)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	for i := range artworks {
		artist, err := s.User(artworks[i].UserID)
		if err != nil {
			return nil, err
		}
		artworks[i].Artist = &artist
	}
	return artworks, nil
}

func (s *Service) Artwork(id int64) (models.Artwork, error) {
	row := s.db.QueryRow(`SELECT id, user_id, title, image_url, caption, created_at FROM artworks WHERE id = ?`, id)
	artwork, err := scanArtwork(row)
	if errors.Is(err, sql.ErrNoRows) {
		return artwork, ErrNotFound
	}
	if err != nil {
		return artwork, err
	}
	artist, err := s.User(artwork.UserID)
	if err != nil {
		return artwork, err
	}
	comments, err := s.Comments(id)
	if err != nil {
		return artwork, err
	}
	artwork.Artist = &artist
	artwork.Comments = comments
	return artwork, nil
}

func (s *Service) CreateComment(ctx context.Context, reviewerID, artworkID int64, input CreateCommentInput) (models.Comment, error) {
	input.Body = strings.TrimSpace(input.Body)
	if input.Body == "" {
		return models.Comment{}, fmt.Errorf("body is required")
	}
	if input.RatingCNZ < 1 || input.RatingCNZ > 5 {
		return models.Comment{}, fmt.Errorf("rating_cnz must be 1, 2, 3, 4, or 5")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Comment{}, err
	}
	defer tx.Rollback()

	var reviewerBalance float64
	err = tx.QueryRow(`SELECT cnz_balance FROM users WHERE id = ?`, reviewerID).Scan(&reviewerBalance)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Comment{}, ErrNotFound
	}
	if err != nil {
		return models.Comment{}, err
	}

	var artistID int64
	err = tx.QueryRow(`SELECT user_id FROM artworks WHERE id = ?`, artworkID).Scan(&artistID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Comment{}, ErrNotFound
	}
	if err != nil {
		return models.Comment{}, err
	}
	if artistID == reviewerID {
		return models.Comment{}, ErrForbidden
	}

	spendAmount := float64(input.RatingCNZ)
	artistAmount := spendAmount * ArtistEarnRate
	if reviewerBalance < spendAmount {
		return models.Comment{}, ErrInsufficientCNZ
	}

	result, err := tx.Exec(`INSERT INTO comments (artwork_id, user_id, rating_cnz, body, artist_cnz) VALUES (?, ?, ?, ?, ?)`,
		artworkID, reviewerID, input.RatingCNZ, input.Body, artistAmount)
	if err != nil {
		return models.Comment{}, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return models.Comment{}, err
	}

	if _, err = tx.Exec(`UPDATE users SET cnz_balance = cnz_balance - ? WHERE id = ?`, spendAmount, reviewerID); err != nil {
		return models.Comment{}, err
	}
	if _, err = tx.Exec(`UPDATE users SET cnz_balance = cnz_balance + ? WHERE id = ?`, artistAmount, artistID); err != nil {
		return models.Comment{}, err
	}
	if _, err = tx.Exec(`INSERT INTO cnz_ledger (user_id, artwork_id, amount, direction, reason) VALUES (?, ?, ?, ?, ?)`,
		reviewerID, artworkID, spendAmount, "debit", "comment_rating"); err != nil {
		return models.Comment{}, err
	}
	if _, err = tx.Exec(`INSERT INTO cnz_ledger (user_id, artwork_id, amount, direction, reason) VALUES (?, ?, ?, ?, ?)`,
		artistID, artworkID, artistAmount, "credit", "artist_earning"); err != nil {
		return models.Comment{}, err
	}

	if err = tx.Commit(); err != nil {
		return models.Comment{}, err
	}
	return s.Comment(commentID)
}

func (s *Service) Comments(artworkID int64) ([]models.Comment, error) {
	rows, err := s.db.Query(`SELECT id, artwork_id, user_id, rating_cnz, body, artist_cnz, created_at FROM comments WHERE artwork_id = ? ORDER BY created_at ASC, id ASC`, artworkID)
	if err != nil {
		return nil, err
	}

	comments := []models.Comment{}
	for rows.Next() {
		comment, err := scanComment(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	for i := range comments {
		author, err := s.User(comments[i].UserID)
		if err != nil {
			return nil, err
		}
		comments[i].Author = &author
	}
	return comments, nil
}

func (s *Service) Comment(id int64) (models.Comment, error) {
	row := s.db.QueryRow(`SELECT id, artwork_id, user_id, rating_cnz, body, artist_cnz, created_at FROM comments WHERE id = ?`, id)
	comment, err := scanComment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return comment, ErrNotFound
	}
	if err != nil {
		return comment, err
	}
	author, err := s.User(comment.UserID)
	if err != nil {
		return comment, err
	}
	comment.Author = &author
	return comment, nil
}

func (s *Service) LedgerEntries() ([]models.LedgerEntry, error) {
	rows, err := s.db.Query(`SELECT id, user_id, artwork_id, amount, direction, reason, created_at FROM cnz_ledger ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LedgerEntry
	for rows.Next() {
		var entry models.LedgerEntry
		var artworkID sql.NullInt64
		if err := rows.Scan(&entry.ID, &entry.UserID, &artworkID, &entry.Amount, &entry.Direction, &entry.Reason, &entry.CreatedAt); err != nil {
			return nil, err
		}
		if artworkID.Valid {
			entry.ArtworkID = &artworkID.Int64
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanArtwork(row rowScanner) (models.Artwork, error) {
	var artwork models.Artwork
	err := row.Scan(&artwork.ID, &artwork.UserID, &artwork.Title, &artwork.ImageURL, &artwork.Caption, &artwork.CreatedAt)
	return artwork, err
}

func scanComment(row rowScanner) (models.Comment, error) {
	var comment models.Comment
	err := row.Scan(&comment.ID, &comment.ArtworkID, &comment.UserID, &comment.RatingCNZ, &comment.Body, &comment.ArtistCNZ, &comment.CreatedAt)
	return comment, err
}
