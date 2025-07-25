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

var Mode string
var Domain string
var WebUrl string

var CorsAllowed string

func Load() {
	if os.Getenv("READ_ENV_FILE") != "skip" {
		err := godotenv.Load()
		assert.AssertErr(err)
	}

	OauthClientId = env("OAUTH_CLIENT_ID", "")
	OauthClientSecret = env("OAUTH_CLIENT_SECRET", "")
	JwtSecret = []byte(env("JWT_SECRET", ""))
	DbConnStr = env("DB_CONN_STR", "")

	Mode = env("MODE", "dev")
	Domain = env("DOMAIN", "localhost:8888")
	WebUrl = env("WEB_URL", "http://localhost:5173")

	CorsAllowed = env("CORS_ALLOWED_ORIGINS", "http://localhost:5173")
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = fallback
	}

	assert.Assertf(len(value) > 0, "expected '%s' to be in environment\r\n", key)
	return value
}
