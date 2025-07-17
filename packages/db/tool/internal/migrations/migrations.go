package migrations

import (
	"database/sql"
	"fmt"
)

type Migration struct {
	Timestamp string
	Label     string
	Up        string
	Down      string
}

func EnsureExists(db *sql.DB) error {
	rows, err := db.Query("select name from sqlite_master where type='table' and name='migrations';")
	exists := err == nil && rows.Next()
	rows.Close()

	if exists {
		return nil
	}

	if _, err := db.Exec("create table migrations (timestamp integer primary key, label text);"); err != nil {
		return fmt.Errorf("could not create migrations table: %w", err)
	}

	return nil

}

func Up(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Could not begin migration transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("insert into migrations (timestamp, label) values (?, ?);")
	if err != nil {
		return fmt.Errorf("Could not create insertion statement: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(migration.Timestamp, migration.Label); err != nil {
		return fmt.Errorf("Could not insert migration '%s': %w", migration.Label, err)
	}

	// Test if up migration is running
	if _, err := tx.Exec(migration.Up); err != nil {
		return fmt.Errorf("Could not execute up migration ('%s'): %w", migration.Label, err)
	}

	// Validate that down migration is also valid given the previous up migration
	if _, err := tx.Exec(migration.Down); err != nil {
		return fmt.Errorf("Could not execute down migration ('%s'): %w", migration.Label, err)
	}

	// If everything was fine, re-execute the up migration and continue
	if _, err := tx.Exec(migration.Up); err != nil {
		return fmt.Errorf("Could not execute up migration ('%s'): %w", migration.Label, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Migration transaction failed: %w", err)
	}

	return nil
}

func Down(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Could not begin rollback transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("delete from migrations where timestamp=? and label=?;")
	if err != nil {
		return fmt.Errorf("Could not create delete statement: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(migration.Timestamp, migration.Label); err != nil {
		return fmt.Errorf("Could not delete migration '%s': %w", migration.Label, err)
	}

	if _, err := tx.Exec(migration.Down); err != nil {
		return fmt.Errorf("Could not rollback migration ('%s'): %w", migration.Label, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Migration rollback failed: %w", err)
	}

	return nil
}
