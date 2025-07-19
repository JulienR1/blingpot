package main

import (
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

	err := server.Run(&config)
	assert.AssertErr(err)
}
