package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RunMigrations(db *sql.DB) error {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	for _, file := range files {

		filename := filepath.Base(file)

		version, err := strconv.Atoi(strings.Split(filename, "_")[0])
		if err != nil {
			continue
		}

		var count int

		err = db.QueryRow(
			"SELECT COUNT(*) FROM schema_migrations WHERE version=?",
			version,
		).Scan(&count)

		if err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		fmt.Println("Applying migration:", filename)

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err = tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return err
		}

		_, err = tx.Exec(
			"INSERT INTO schema_migrations(version) VALUES(?)",
			version,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
