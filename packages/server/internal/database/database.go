package database

import (
	"database/sql"
	"fmt"

	"github.com/julienr1/blingpot/internal/env"
)

type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
}

func Open() (*sql.DB, error) {
	return sql.Open("sqlite3", env.DbConnStr)
}

func Transaction(callback func(tx Querier) error) error {
	db, err := Open()
	if err != nil {
		return fmt.Errorf("database.Transaction: could not open db, %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("database.Transaction: could not begin transaction, %w", err)
	}
	defer tx.Rollback()

	if err = callback(tx); err != nil {
		return fmt.Errorf("database.Transaction: could not execute transaction body, %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("database.Transaction: could not commit transaction, %w", err)
	}

	return nil
}
