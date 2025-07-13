package server

import (
	"log"
	"net/http"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/auth"
)

func Run(config *ServerConfig) error {
	assert.Assert(config != nil, "server config is nil")

	auth := auth.New(config.ServerUrl(), config.WebUrl)
	http.HandleFunc("GET /oauth2/authenticate", auth.HandleAuth)
	http.HandleFunc("GET /oauth2/callback", auth.HandleAuthCallback)
	http.HandleFunc("POST /oauth2/revoke", auth.HandleRevoke)

	http.Handle("/*", http.RedirectHandler("/", http.StatusTemporaryRedirect))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	log.Println("Listening on", config.ServerUrl())
	return http.ListenAndServe(config.Endpoint(), nil)
}
