package category

import (
	"errors"
	"fmt"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
)

type Category struct {
	Id         int
	Label      string
	Foreground string
	Background string
	IconName   string
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
	err = row.Scan(&c.Id, &c.Label, &c.Foreground, &c.Background, &c.IconName)
	if err != nil {
		return nil, fmt.Errorf("category.FindById: %w %w", ErrCategoryNotFound, err)
	}

	return &c, nil
}

func FindAll(db database.Querier) ([]Category, error) {
	stmt, err := db.Prepare("select id, label, color_fg, color_bg, icon_name from categories;")
	assert.AssertErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []Category{}, fmt.Errorf("category.FindAll: %w", err)
	}

	var categories []Category
	for rows.Next() {
		var c Category
		if err = rows.Scan(&c.Id, &c.Label, &c.Foreground, &c.Background, &c.IconName); err != nil {
			return []Category{}, fmt.Errorf("category.FindAll: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (c Category) Dto() dtos.Category {
	return dtos.Category{
		Id:       c.Id,
		Label:    c.Label,
		IconName: c.IconName,
		Color: dtos.Color{
			Foregound:  c.Foreground,
			Background: c.Background,
		},
	}
}
