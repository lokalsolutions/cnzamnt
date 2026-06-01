package seed

import (
	"database/sql"
	"fmt"
)

func Ensure(database *sql.DB) error {
	var count int
	if err := database.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT INTO users (display_name, handle, cnz_balance) VALUES (?, ?, ?)`, "Demo Artist", "demo_artist", 5000)
	if err != nil {
		return err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO cnz_ledger (user_id, amount, direction, reason) VALUES (?, ?, ?, ?)`, userID, 5000, "credit", "signup_bonus"); err != nil {
		return err
	}

	return tx.Commit()
}
