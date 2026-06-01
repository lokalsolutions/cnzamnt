package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"cnzamnt/backend/internal/db"
	"cnzamnt/backend/internal/models"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	database, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	t.Cleanup(func() {
		_ = database.Close()
	})
	return New(database)
}

func TestCreateUserStartsWith5000CNZ(t *testing.T) {
	svc := newTestService(t)

	user, err := svc.CreateUser(CreateUserInput{
		DisplayName: "Mira",
		Handle:      "mira",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if user.CNZBalance != StartingCNZ {
		t.Fatalf("expected %v CNZ, got %v", StartingCNZ, user.CNZBalance)
	}

	entries, err := svc.LedgerEntries()
	if err != nil {
		t.Fatalf("load ledger: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 ledger entry, got %d", len(entries))
	}
	if entries[0].UserID != user.ID || entries[0].Amount != StartingCNZ || entries[0].Direction != "credit" || entries[0].Reason != "signup_bonus" {
		t.Fatalf("unexpected signup ledger entry: %+v", entries[0])
	}
}

func TestCommentSpendsCNZAndCreditsArtist(t *testing.T) {
	svc := newTestService(t)
	artist, reviewer, artwork := createArtistReviewerArtwork(t, svc)

	comment, err := svc.CreateComment(context.Background(), reviewer.ID, artwork.ID, CreateCommentInput{
		RatingCNZ: 5,
		Body:      "Strong composition and color balance.",
	})
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}
	if comment.RatingCNZ != 5 {
		t.Fatalf("expected rating 5, got %d", comment.RatingCNZ)
	}

	updatedReviewer, err := svc.User(reviewer.ID)
	if err != nil {
		t.Fatalf("load reviewer: %v", err)
	}
	if updatedReviewer.CNZBalance != 4995 {
		t.Fatalf("expected reviewer balance 4995, got %v", updatedReviewer.CNZBalance)
	}

	updatedArtist, err := svc.User(artist.ID)
	if err != nil {
		t.Fatalf("load artist: %v", err)
	}
	if updatedArtist.CNZBalance != 5000.5 {
		t.Fatalf("expected artist balance 5000.5, got %v", updatedArtist.CNZBalance)
	}

	comments, err := svc.Comments(artwork.ID)
	if err != nil {
		t.Fatalf("load comments: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].Author == nil || comments[0].Author.ID != reviewer.ID {
		t.Fatalf("expected comment author to be reviewer, got %+v", comments[0].Author)
	}
}

func TestCommentRejectsInsufficientCNZ(t *testing.T) {
	svc := newTestService(t)
	_, reviewer, artwork := createArtistReviewerArtwork(t, svc)

	if _, err := svc.db.Exec(`UPDATE users SET cnz_balance = 3 WHERE id = ?`, reviewer.ID); err != nil {
		t.Fatalf("set reviewer balance: %v", err)
	}

	_, err := svc.CreateComment(context.Background(), reviewer.ID, artwork.ID, CreateCommentInput{
		RatingCNZ: 5,
		Body:      "I like the line work.",
	})
	if !errors.Is(err, ErrInsufficientCNZ) {
		t.Fatalf("expected ErrInsufficientCNZ, got %v", err)
	}

	comments, err := svc.Comments(artwork.ID)
	if err != nil {
		t.Fatalf("load comments: %v", err)
	}
	if len(comments) != 0 {
		t.Fatalf("expected no comments, got %d", len(comments))
	}

	entries, err := svc.LedgerEntries()
	if err != nil {
		t.Fatalf("load ledger: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected only signup ledger entries, got %d", len(entries))
	}
	for _, entry := range entries {
		if entry.Reason != "signup_bonus" {
			t.Fatalf("expected no comment ledger entries, got %+v", entry)
		}
	}
}

func TestCommentCreatesLedgerEntries(t *testing.T) {
	svc := newTestService(t)
	artist, reviewer, artwork := createArtistReviewerArtwork(t, svc)

	_, err := svc.CreateComment(context.Background(), reviewer.ID, artwork.ID, CreateCommentInput{
		RatingCNZ: 4,
		Body:      "The lighting makes the subject read clearly.",
	})
	if err != nil {
		t.Fatalf("create comment: %v", err)
	}

	entries, err := svc.LedgerEntries()
	if err != nil {
		t.Fatalf("load ledger: %v", err)
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 ledger entries, got %d", len(entries))
	}

	debit := entries[2]
	if debit.UserID != reviewer.ID || debit.ArtworkID == nil || *debit.ArtworkID != artwork.ID || debit.Amount != 4 || debit.Direction != "debit" || debit.Reason != "comment_rating" {
		t.Fatalf("unexpected debit entry: %+v", debit)
	}

	credit := entries[3]
	if credit.UserID != artist.ID || credit.ArtworkID == nil || *credit.ArtworkID != artwork.ID || credit.Amount != 0.4 || credit.Direction != "credit" || credit.Reason != "artist_earning" {
		t.Fatalf("unexpected credit entry: %+v", credit)
	}
}

func TestCommentRatingMustBeOneThroughFive(t *testing.T) {
	svc := newTestService(t)
	_, reviewer, artwork := createArtistReviewerArtwork(t, svc)

	_, err := svc.CreateComment(context.Background(), reviewer.ID, artwork.ID, CreateCommentInput{
		RatingCNZ: 6,
		Body:      "Too much.",
	})
	if err == nil {
		t.Fatal("expected rating validation error")
	}
}

func createArtistReviewerArtwork(t *testing.T, svc *Service) (artist, reviewer models.User, artwork models.Artwork) {
	t.Helper()

	createdArtist, err := svc.CreateUser(CreateUserInput{DisplayName: "Ari", Handle: "ari"})
	if err != nil {
		t.Fatalf("create artist: %v", err)
	}
	createdReviewer, err := svc.CreateUser(CreateUserInput{DisplayName: "Nia", Handle: "nia"})
	if err != nil {
		t.Fatalf("create reviewer: %v", err)
	}
	createdArtwork, err := svc.CreateArtwork(createdArtist.ID, CreateArtworkInput{
		Title:    "Blue Room",
		ImageURL: "https://example.test/blue-room.jpg",
		Caption:  "Study in blue.",
	})
	if err != nil {
		t.Fatalf("create artwork: %v", err)
	}
	return createdArtist, createdReviewer, createdArtwork
}
