package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/julienr1/blingpot/internal/assert"
)

var OauthClientId string
var OauthClientSecret string
var JwtSecret []byte
var DbConnStr string

func Load() {
	err := godotenv.Load()
	assert.AssertErr(err)

	OauthClientId = env("OAUTH_CLIENT_ID")
	OauthClientSecret = env("OAUTH_CLIENT_SECRET")
	JwtSecret = []byte(env("JWT_SECRET"))
	DbConnStr = env("DB_CONN_STR")
}

func env(key string) string {
	value := os.Getenv(key)
	assert.Assertf(len(value) > 0, "expected '%s' to be in environment", key)
	return value
}
