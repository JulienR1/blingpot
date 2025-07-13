package main

import (
	"database/sql"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/env"
	"github.com/julienr1/blingpot/internal/server"
	_ "github.com/mattn/go-sqlite3"
)

var config = server.ServerConfig{
	Protocol: server.HTTP,
	Domain:   "localhost",
	Port:     8888,
	WebUrl:   "http://localhost:5173",
}

func main() {
	env.Load()

	db, err := sql.Open("sqlite3", env.DbConnStr)
	assert.AssertErr(err)
	db.Exec("create table foo ( id integer not null primary key, name text );")
	db.Close()

	err = server.Run(&config)
	assert.AssertErr(err)
}
