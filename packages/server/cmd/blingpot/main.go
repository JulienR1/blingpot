package main

import (
	"github.com/joho/godotenv"
	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/server"
)

var config = server.ServerConfig{
	Protocol: server.HTTP,
	Domain:   "localhost",
	Port:     8888,
	WebUrl:   "http://localhost:5173",
}

func main() {
	err := godotenv.Load()
	assert.AssertErr(err)

	err = server.Run(&config)
	assert.AssertErr(err)
}
