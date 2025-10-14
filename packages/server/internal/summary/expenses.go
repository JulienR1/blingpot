package summary

import (
	"time"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
)

func GetTotal(db database.Querier, start, end time.Time) (total int) {
	assert.Assert(start.Unix() < end.Unix(), "summary.GetTotal: invalid time interval")

	stmt, err := db.Prepare(`
		select ifnull(sum(amount), 0) as subtotal
		from expenses
		where datetime >= ? and datetime <= ?;
	`)
	assert.AssertErr(err)
	defer stmt.Close()

	row := stmt.QueryRow(start.Unix(), end.Unix())
	err = row.Scan(&total)
	assert.AssertErr(err)

	return total
}

func GetCategoryTotals(db database.Querier, start, end time.Time) map[int]int {
	assert.Assert(start.Unix() < end.Unix(), "summary.GetCategoryTotals: invalid time interval")

	stmt, err := db.Prepare(`
		select category_id, sum(amount) as subtotal
		from (
			select category_id, amount
			from expenses
			where datetime >= ? and datetime <= ?
			union all
			select id as category_id, 0 as amount
			from categories
		)
		group by category_id;
	`)
	assert.AssertErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(start.Unix(), end.Unix())
	assert.AssertErr(err)
	defer rows.Close()

	var subtotals = make(map[int]int)
	for rows.Next() {
		var categoryId, subtotal int
		rows.Scan(&categoryId, &subtotal)
		subtotals[categoryId] = subtotal
	}

	return subtotals
}
