package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/env"
	"github.com/julienr1/blingpot/internal/profile"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	verifierId uint8
	verifiers  map[uint8]string

	config *oauth2.Config
}

func New(url string) *Auth {
	return &Auth{
		verifierId: 0,
		verifiers:  make(map[uint8]string),

		config: &oauth2.Config{
			ClientID:     env.OauthClientId,
			ClientSecret: env.OauthClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email ",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint:    google.Endpoint,
			RedirectURL: fmt.Sprintf("%s/oauth2/callback", url),
		},
	}
}

func (a *Auth) HandleAuth(w http.ResponseWriter, r *http.Request) {
	if p, ok := r.Context().Value("profile").(profile.Profile); ok {
		fmt.Println("Found connected profile:", p)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	a.verifiers[a.verifierId] = oauth2.GenerateVerifier()
	state := strconv.Itoa(int(a.verifierId))
	challenge := oauth2.S256ChallengeOption(a.verifiers[a.verifierId])
	a.verifierId++

	url := a.config.AuthCodeURL(state, oauth2.AccessTypeOffline, challenge)
	http.Redirect(w, r, url, http.StatusFound)
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

	var userInfo dtos.OauthUserInfo
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

	if err = database.Transaction(func(tx database.Querier) error {
		return profile.StoreProfile(tx, userInfo.Sub, userInfo.GivenName, userInfo.FamilyName, userInfo.Email, userInfo.Picture, token)
	}); err != nil {
		log.Println("HandleAuthCallback: could not create profile,", err)
		http.Error(w, "could not create profile", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    jwt,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (a *Auth) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	assert.AssertErr(err)

	cookie.Path = "/"
	cookie.Expires = time.Now().Add(5 * time.Minute)
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (a *Auth) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	p, ok := r.Context().Value("profile").(profile.Profile)
	assert.Assert(ok, "HandleRevoke: profile should have been set in context")

	db, err := database.Open()
	assert.AssertErr(err)
	defer db.Close()

	token, err := p.ProviderToken.Value()
	if err != nil {
		http.Error(w, "could not find provider token for profile", http.StatusInternalServerError)
		return
	}

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

	if err = profile.ClearProviderToken(db, p.Sub); err != nil {
		http.Error(w, "local token was not revoked", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
}
