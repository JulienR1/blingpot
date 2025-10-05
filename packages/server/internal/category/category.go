package category

import (
	"errors"
	"fmt"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
)

type Color struct {
	Background string
	Foregound  string
}

type Category struct {
	Id       int
	Label    string
	Color    Color
	IconName string
}

var ErrCategoryNotFound = errors.New("Could not find category")

func FindById(db database.Querier, id int) (*Category, error) {
	var c Category

	query := `
		select id, label, color_fg, color_bg, icon_name
		from categories
		where id = ?;
	`

	stmt, err := db.Prepare(query)
	assert.AssertErr(err)
	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(&c.Id, &c.Label, &c.Color.Foregound, &c.Color.Background, &c.IconName)
	if err != nil {
		return nil, fmt.Errorf("category.FindById: %w %w", ErrCategoryNotFound, err)
	}

	return &c, nil
}
