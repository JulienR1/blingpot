package profile

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
)

type Profile struct {
	Sub       string
	FirstName string
	LastName  string
	Email     string
	Picture   sql.NullString
}

var ProfileNotFound = errors.New("Could not find profile")

func FindBySub(db database.Querier, sub string) (*Profile, error) {
	var p Profile

	stmt, err := db.Prepare("select sub, first_name, last_name, email, picture from profiles where sub=?;")
	assert.AssertErr(err)
	defer stmt.Close()

	row := stmt.QueryRow(sub)
	err = row.Scan(&p.Sub, &p.FirstName, &p.LastName, &p.Email, &p.Picture)
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

