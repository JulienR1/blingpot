package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienr1/blingpot/internal/assert"
)

func jwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	assert.Assert(len(secret) > 0, "could not find jwt secret in environment")
	return []byte(secret)
}

func CreateToken(user *UserInfo) (string, error) {
	exp := time.Now().Add(time.Minute * 10)
	claims := jwt.MapClaims{"sub": user.Sub, "exp": exp.Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
}

func VerifyToken(token string) bool {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return jwtSecret(), nil
	})
	return err == nil && t.Valid
}
