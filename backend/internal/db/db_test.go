package db

import (
	"path/filepath"
	"testing"
)

func TestOpenCreatesSQLiteDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cnzamnt.db")

	database, err := Open(path)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	var tableName string
	err = database.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'schema_migrations'`).Scan(&tableName)
	if err != nil {
		t.Fatalf("expected schema_migrations table: %v", err)
	}
	if tableName != "schema_migrations" {
		t.Fatalf("expected schema_migrations table, got %q", tableName)
	}
}
