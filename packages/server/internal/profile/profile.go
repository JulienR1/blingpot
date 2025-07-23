package profile

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"golang.org/x/oauth2"
)

type Profile struct {
	Sub           string
	FirstName     string
	LastName      string
	Email         string
	Picture       sql.NullString
	ProviderToken sql.NullString
}

var ProfileNotFound = errors.New("Could not find profile")

func FindBySub(db database.Querier, sub string) (*Profile, error) {
	var p Profile

	query := `
        select p.sub, first_name, last_name, email, picture, access_token
        from profiles p
        left join (
            select sub, access_token
            from provider_tokens
        ) pt on pt.sub = p.sub
        where p.sub = ?;
    `

	stmt, err := db.Prepare(query)
	assert.AssertErr(err)
	defer stmt.Close()

	row := stmt.QueryRow(sub)
	err = row.Scan(&p.Sub, &p.FirstName, &p.LastName, &p.Email, &p.Picture, &p.ProviderToken)
	if err != nil {
		return nil, fmt.Errorf("profile.FindBySub: %w %w", ProfileNotFound, err)
	}

	return &p, nil
}

func Create(db database.Querier, sub, firstName, lastName, email, picture string) error {
	stmt, err := db.Prepare("insert into profiles (sub, first_name, last_name, email, picture, provider_id) values (?, ?, ?, ?, ?, ?);")
	assert.AssertErr(err)
	defer stmt.Close()

	var pic = sql.NullString{Valid: false}
	if len(picture) > 0 {
		pic = sql.NullString{String: picture, Valid: true}
	}

	_, err = stmt.Exec(sub, firstName, lastName, email, pic, "google")
	if err != nil {
		return fmt.Errorf("profile.Create: %w", err)
	}

	return nil
}

func StoreProfile(db database.Querier, sub, firstName, lastName, email, picture string, token *oauth2.Token) error {
	query := `
        insert into profiles (sub, first_name, last_name, email, picture, provider_id) values (?, ?, ?, ?, ?, ?)
        on conflict (sub) do
        update set first_name = excluded.first_name, last_name = excluded.last_name, picture = excluded.picture;
    `
	stmt, err := db.Prepare(query)
	assert.AssertErr(err)

	var pic = sql.NullString{Valid: false}
	if len(picture) > 0 {
		pic = sql.NullString{String: picture, Valid: true}
	}

	_, err = stmt.Exec(sub, firstName, lastName, email, pic, "google")
	stmt.Close()
	if err != nil {
		return fmt.Errorf("profile.StoreProfile: %w", err)
	}

	query = `
       insert into provider_tokens (sub, access_token, refresh_token) values (?, ?, ?)
       on conflict (sub) do
       update set access_token = excluded.access_token, refresh_token = excluded.refresh_token;
  `
	stmt, err = db.Prepare(query)
	assert.AssertErr(err)

	_, err = stmt.Exec(sub, token.AccessToken, token.RefreshToken)
	stmt.Close()
	if err != nil {
		return fmt.Errorf("profile.StoreProfile: %w", err)
	}

	return nil
}

func ClearProviderToken(db database.Querier, sub string) error {
	stmt, err := db.Prepare("delete from provider_tokens where sub = ?;")
	assert.AssertErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(sub)
	if err != nil {
		return fmt.Errorf("profile.ClearProviderToken: %w", err)
	}

	return nil
}
