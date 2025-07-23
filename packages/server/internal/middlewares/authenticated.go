package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/auth"
	"github.com/julienr1/blingpot/internal/database"
)

type AuthenticationMode int

const (
	Optional AuthenticationMode = iota
	Required
)

func Authenticate(next http.HandlerFunc, mode AuthenticationMode) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")
		if errors.Is(err, http.ErrNoCookie) && mode == Optional {
			next(w, r)
			return
		}
		if err != nil {
			http.Error(w, "no authentication token", http.StatusUnauthorized)
			return
		}

		if auth.VerifyToken(cookie.Value) == false && mode == Optional {
			next(w, r)
			return
		}

		db, err := database.Open()
		assert.AssertErr(err)
		defer db.Close()

		p, err := auth.FindProfileFromToken(db, cookie.Value)
		if err != nil || p == nil {
			http.Error(w, "could not find profile", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "profile", *p)
		next(w, r.WithContext(ctx))
	})
}

func Authenticated(next http.HandlerFunc) http.Handler {
	return Authenticate(next, Required)
}
