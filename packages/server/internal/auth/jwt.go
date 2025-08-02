package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/env"
	"github.com/julienr1/blingpot/internal/profile"
)

var InvalidTokenErr = errors.New("Could not parse or verify token")
var MalformedTokenErr = errors.New("Could not find field on token")

func CreateToken(user *dtos.OauthUserInfo) (string, error) {
	exp := jwt.NewNumericDate(time.Now().Add(time.Minute * 10))
	claims := jwt.RegisteredClaims{Subject: user.Sub, ExpiresAt: exp}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(env.JwtSecret)
}

func VerifyToken(token string) bool {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return env.JwtSecret, nil
	})

	return err == nil && t.Valid
}

func FindProfileFromToken(db database.Querier, token string) (*profile.Profile, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return env.JwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("auth.FindProfileByToken: %w", InvalidTokenErr)
	}

	if !t.Valid {
		return nil, fmt.Errorf("auth.FindProfileByToken: %w", InvalidTokenErr)
	}

	sub, err := t.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("auth.FindProfileByToken: %w", MalformedTokenErr)
	}

	return profile.FindBySub(db, sub)
}
