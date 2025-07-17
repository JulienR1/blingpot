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

type SQLQuerier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
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

func Up(db SQLQuerier, migration Migration) error {
	stmt, err := db.Prepare("insert into migrations (timestamp, label) values (?, ?);")
	if err != nil {
		return fmt.Errorf("Could not create insertion statement: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(migration.Timestamp, migration.Label); err != nil {
		return fmt.Errorf("Could not insert migration '%s': %w", migration.Label, err)
	}

	// Test if up migration is running
	if _, err := db.Exec(migration.Up); err != nil {
		return fmt.Errorf("Could not execute up migration ('%s'): %w", migration.Label, err)
	}

	// Validate that down migration is also valid given the previous up migration
	if _, err := db.Exec(migration.Down); err != nil {
		return fmt.Errorf("Could not execute down migration ('%s'): %w", migration.Label, err)
	}

	// If everything was fine, re-execute the up migration and continue
	if _, err := db.Exec(migration.Up); err != nil {
		return fmt.Errorf("Could not execute up migration ('%s'): %w", migration.Label, err)
	}

	return nil
}

func Down(db SQLQuerier, migration Migration) error {
	stmt, err := db.Prepare("delete from migrations where timestamp=? and label=?;")
	if err != nil {
		return fmt.Errorf("Could not create delete statement: %w", err)
	}
	defer stmt.Close()

	if _, err = stmt.Exec(migration.Timestamp, migration.Label); err != nil {
		return fmt.Errorf("Could not delete migration '%s': %w", migration.Label, err)
	}

	if _, err := db.Exec(migration.Down); err != nil {
		return fmt.Errorf("Could not rollback migration ('%s'): %w", migration.Label, err)
	}

	return nil
}
