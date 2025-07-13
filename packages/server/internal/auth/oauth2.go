package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienr1/blingpot/internal/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Auth struct {
	verifierId uint8
	verifiers  map[uint8]string
	appUrl     string

	config *oauth2.Config
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func New(serverUrl, appUrl string) *Auth {
	clientId := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")

	assert.Assert(len(clientId) > 0, "missing oauth client id in environment")
	assert.Assert(len(clientSecret) > 0, "missing oauth client secret in environment")

	return &Auth{
		verifierId: 0,
		verifiers:  make(map[uint8]string),
		appUrl:     appUrl,

		config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
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

	fmt.Println(userInfo, jwt)

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
