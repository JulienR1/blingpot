package main

import (
	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/env"
	"github.com/julienr1/blingpot/internal/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	env.Load()

	var config = server.ServerConfig{
		Protocol: server.HTTPS,
		Domain:   env.Domain,
		Port:     8888,
		WebDir:   env.WebDir,
	}

	if env.Mode == "dev" {
		config.Protocol = server.HTTP
	}

	err := server.Run(&config)
	assert.AssertErr(err)
}
