package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email ",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint:    google.Endpoint,
		RedirectURL: "http://localhost:8888/oauth2/callback",
	}

	var verifierId uint8 = 0
	var verifiers = make(map[uint8]string)

	http.HandleFunc("GET /oauth2/authenticate", func(w http.ResponseWriter, r *http.Request) {
		verifiers[verifierId] = oauth2.GenerateVerifier()
		url := config.AuthCodeURL(strconv.Itoa(int(verifierId)), oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifiers[verifierId]))
		verifierId++

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
	http.HandleFunc("GET /oauth2/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		verifierIdStr := q.Get("state")
		vId, _ := strconv.Atoi(verifierIdStr)
		v := verifiers[uint8(vId)]
		delete(verifiers, uint8(vId))

		code := q.Get("code")
		ctx := context.Background()
		token, err := config.Exchange(ctx, code, oauth2.VerifierOption(v))
		if err != nil {
			http.Error(w, "could not authenticate", http.StatusUnauthorized)
			return
		}

		client := config.Client(ctx, token)
		rr, _ := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
		b, _ := io.ReadAll(rr.Body)
		fmt.Println(w, string(b))

		http.Redirect(w, r, "http://localhost:5173", http.StatusTemporaryRedirect)
	})
	http.HandleFunc("POST /oauth2/revoke", func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token[:len("Bearer ")] == "Bearer " {
			token = token[len("Bearer "):]
			url := fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token)
			r, _ := http.NewRequest("POST", url, nil)

			c := http.Client{}
			res, err := c.Do(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer res.Body.Close()
		}

		w.WriteHeader(http.StatusOK)
	})

	http.Handle("/*", http.RedirectHandler("/", http.StatusTemporaryRedirect))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	fmt.Println("Listening on http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
