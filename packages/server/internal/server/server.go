package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/auth"
	"github.com/julienr1/blingpot/internal/category"
	"github.com/julienr1/blingpot/internal/env"
	"github.com/julienr1/blingpot/internal/expense"
	"github.com/julienr1/blingpot/internal/middlewares"
	"github.com/julienr1/blingpot/internal/profile"
	"github.com/julienr1/blingpot/internal/summary"
)

func rootHandler(webdir string) func(http.ResponseWriter, *http.Request) {
	if env.Mode == "dev" {
		url, _ := url.Parse(env.AppServerUrl)
		proxy := httputil.NewSingleHostReverseProxy(url)
		return func(w http.ResponseWriter, r *http.Request) { proxy.ServeHTTP(w, r) }
	}
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/index.html", webdir))
	}
}

func Run(config *ServerConfig) error {
	assert.Assert(config != nil, "server config is nil")

	auth := auth.New(config.ServerUrl())
	http.Handle("GET /oauth2/authenticate", middlewares.Authenticate(auth.HandleAuth, middlewares.Optional))
	http.HandleFunc("GET /oauth2/callback", auth.HandleAuthCallback)
	http.Handle("POST /oauth2/refresh", middlewares.Authenticated(auth.HandleRefresh))
	http.Handle("POST /oauth2/revoke", middlewares.Authenticated(auth.HandleRevoke))

	http.Handle("GET /profiles", middlewares.Authenticated(profile.HandleFindAll))
	http.Handle("GET /profiles/me", middlewares.Authenticated(profile.HandleFindMe))

	http.Handle("GET /categories", middlewares.Authenticated(category.HandleFindAll))

	http.Handle("GET /expenses", middlewares.Authenticated(expense.HandleFind))
	http.Handle("POST /expenses", middlewares.Authenticated(expense.HandleCreate))

	http.Handle("GET /summary/expenses", middlewares.Authenticated(summary.HandleExpenses))

	fs := http.FileServer(http.Dir(fmt.Sprintf("%s/assets", config.WebDir)))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/*", http.RedirectHandler("/", http.StatusFound))
	http.HandleFunc("/", rootHandler(config.WebDir))

	log.Println("Listening on", config.ServerUrl())
	return http.ListenAndServe(config.Endpoint(), nil)
}
