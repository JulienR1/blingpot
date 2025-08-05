package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/auth"
	"github.com/julienr1/blingpot/internal/middlewares"
	"github.com/julienr1/blingpot/internal/profile"
)

func Run(config *ServerConfig) error {
	assert.Assert(config != nil, "server config is nil")

	auth := auth.New(config.ServerUrl())
	http.Handle("GET /oauth2/authenticate", middlewares.Authenticate(auth.HandleAuth, middlewares.Optional))
	http.HandleFunc("GET /oauth2/callback", auth.HandleAuthCallback)
	http.Handle("POST /oauth2/refresh", middlewares.Authenticated(auth.HandleRefresh))
	http.Handle("POST /oauth2/revoke", middlewares.Authenticated(auth.HandleRevoke))

	http.Handle("GET /profiles", middlewares.Authenticated(profile.HandleFindAll))
	http.Handle("GET /profiles/me", middlewares.Authenticated(profile.HandleFindMe))

	fs := http.FileServer(http.Dir(fmt.Sprintf("%s/assets", config.WebDir)))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/*", http.RedirectHandler("/", http.StatusFound))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", config.WebDir))
	})

	log.Println("Listening on", config.ServerUrl())
	return http.ListenAndServe(config.Endpoint(), nil)
}
