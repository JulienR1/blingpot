package database

import (
	"database/sql"

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
