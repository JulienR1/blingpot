package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/env"
	"github.com/julienr1/blingpot/internal/profile"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	verifierId uint8
	verifiers  map[uint8]string
	appUrl     string

	config *oauth2.Config
}

func New(serverUrl, appUrl string) *Auth {
	return &Auth{
		verifierId: 0,
		verifiers:  make(map[uint8]string),
		appUrl:     appUrl,

		config: &oauth2.Config{
			ClientID:     env.OauthClientId,
			ClientSecret: env.OauthClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email ",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint:    google.Endpoint,
			RedirectURL: fmt.Sprintf("%s/oauth2/callback", serverUrl),
		},
	}
}

func (a *Auth) HandleAuth(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("Authorization"); err == nil {
		token := cookie.Value
		if VerifyToken(token) {
			db, err := database.Open()
			assert.AssertErr(err)
			defer db.Close()

			if profile, err := FindProfileFromToken(db, token); err == nil {
				fmt.Println("Found connected profile:", profile, err)
				http.Redirect(w, r, a.appUrl, http.StatusTemporaryRedirect)
				return
			}
		}
	}

	a.verifiers[a.verifierId] = oauth2.GenerateVerifier()
	state := strconv.Itoa(int(a.verifierId))
	challenge := oauth2.S256ChallengeOption(a.verifiers[a.verifierId])
	a.verifierId++

	url := a.config.AuthCodeURL(state, oauth2.AccessTypeOffline, challenge)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (a *Auth) HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	verifierId, err := strconv.Atoi(r.URL.Query().Get("state"))
	if err != nil {
		http.Error(w, "could not verify authentication request", http.StatusUnauthorized)
		return
	}

	verifier := a.verifiers[uint8(verifierId)]
	delete(a.verifiers, uint8(verifierId))

	ctx := context.Background()
	code := r.URL.Query().Get("code")
	token, err := a.config.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		http.Error(w, "could not complete authentication", http.StatusUnauthorized)
		return
	}

	client := a.config.Client(ctx, token)
	request, err := client.Get("https://openidconnect.googleapis.com/v1/userinfo")
	if err != nil {
		http.Error(w, "could not find user info", http.StatusInternalServerError)
		return
	}

	var userInfo UserInfo
	err = json.NewDecoder(request.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, "could not parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	jwt, err := CreateToken(&userInfo)
	if err != nil {
		http.Error(w, "could not create jwt token:"+err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := database.Open()
	if err != nil {
		http.Error(w, "could not open database:"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = profile.FindBySub(db, userInfo.Sub)
	if errors.Is(err, profile.ProfileNotFound) {
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "fail", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		if err := profile.Create(tx, userInfo.Sub, userInfo.GivenName, userInfo.FamilyName, userInfo.Email, userInfo.Picture); err != nil {
			log.Println("HandleAuthCallback: could not create profile,", err)
			http.Error(w, "could not create profile", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    jwt,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
	})

	http.Redirect(w, r, a.appUrl, http.StatusTemporaryRedirect)
}

func (a *Auth) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	var token = ""
	url := fmt.Sprintf("https://oauth2.googleapis.com/revoke?token=%s", token)

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		http.Error(w, "could not revoke token", http.StatusInternalServerError)
		return
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		http.Error(w, "token was not revoked", response.StatusCode)
		return
	}
	response.Body.Close()
}
