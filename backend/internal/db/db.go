package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	if path == "" {
		path = "data/cnzamnt.db"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	database.SetMaxOpenConns(1)

	if _, err := database.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		database.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := migrate(database); err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

func migrate(database *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS schema_migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			display_name TEXT NOT NULL,
			handle TEXT NOT NULL UNIQUE,
			cnz_balance REAL NOT NULL DEFAULT 5000,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS artworks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			image_url TEXT NOT NULL,
			caption TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			artwork_id INTEGER NOT NULL REFERENCES artworks(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			rating_cnz INTEGER NOT NULL CHECK (rating_cnz IN (1, 2, 3, 4, 5)),
			body TEXT NOT NULL,
			artist_cnz REAL NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS cnz_ledger (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			artwork_id INTEGER REFERENCES artworks(id) ON DELETE SET NULL,
			amount REAL NOT NULL,
			direction TEXT NOT NULL CHECK (direction IN ('debit', 'credit')),
			reason TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_artworks_user_id ON artworks(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_comments_artwork_id ON comments(artwork_id);`,
		`CREATE INDEX IF NOT EXISTS idx_cnz_ledger_user_id ON cnz_ledger(user_id);`,
	}

	for _, statement := range statements {
		if _, err := database.Exec(statement); err != nil {
			return fmt.Errorf("migrate database: %w", err)
		}
	}
	return nil
}
