package server

import (
	"log"
	"net/http"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/auth"
	"github.com/rs/cors"
)

func Run(config *ServerConfig) error {
	assert.Assert(config != nil, "server config is nil")

	mux := http.NewServeMux()

	auth := auth.New(config.ServerUrl(), config.WebUrl)
	mux.HandleFunc("GET /oauth2/authenticate", auth.HandleAuth)
	mux.HandleFunc("GET /oauth2/callback", auth.HandleAuthCallback)
	mux.HandleFunc("POST /oauth2/revoke", auth.HandleRevoke)

	mux.Handle("/*", http.RedirectHandler("/", http.StatusTemporaryRedirect))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{config.WebUrl},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}).Handler(mux)

	log.Println("Listening on", config.ServerUrl())
	return http.ListenAndServe(config.Endpoint(), handler)
}
