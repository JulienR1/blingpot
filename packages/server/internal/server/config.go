package server

import "fmt"

type Protocol string

const (
	HTTP  Protocol = "http"
	HTTPS          = "https"
)

type ServerConfig struct {
	Protocol Protocol
	Domain   string
	Port     int

	WebUrl string
}

func (config *ServerConfig) Endpoint() string {
	return fmt.Sprintf(":%d", config.Port)
}

func (config *ServerConfig) ServerUrl() string {
	return fmt.Sprintf("%s://%s", config.Protocol, config.Domain)
}
