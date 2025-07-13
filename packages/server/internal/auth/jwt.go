package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienr1/blingpot/internal/env"
)

func CreateToken(user *UserInfo) (string, error) {
	exp := time.Now().Add(time.Minute * 10)
	claims := jwt.MapClaims{"sub": user.Sub, "exp": exp.Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(env.JwtSecret)
}

func VerifyToken(token string) bool {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return env.JwtSecret, nil
	})
	return err == nil && t.Valid
}
