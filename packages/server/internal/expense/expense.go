package expense

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/category"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/profile"
)

var InvalidExpenseAmountErr = errors.New("invalid expense amount")
var InvalidExpenseLabelErr = errors.New("invalid expense label")
var ExpenseCreateErr = errors.New("could not create expense")
var ExpenseNotFoundErr = errors.New("could not find expenses")

type Expense struct {
	Id         int
	SpenderId  string
	Label      string
	Amount     int
	Timestamp  time.Time
	AuthorId   string
	CategoryId int
}

func Find(db database.Querier, start, end time.Time) ([]Expense, error) {
	stmt, err := db.Prepare(`
		select id, spender_id, label, amount, datetime, author_id, category_id
		from expenses
		where datetime >= ? and datetime <= ?
		order by datetime desc;
	`)
	assert.AssertErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(start.Unix(), end.Unix())
	if err != nil {
		return []Expense{}, ExpenseNotFoundErr
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var timestamp int64
		var expense Expense
		rows.Scan(&expense.Id, &expense.SpenderId, &expense.Label, &expense.Amount, &timestamp, &expense.AuthorId, &expense.CategoryId)
		expense.Timestamp = time.Unix(timestamp, 0)
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func Create(db database.Querier, label string, amount int, timestamp time.Time, spender, author *profile.Profile, category *category.Category) (id int, err error) {
	assert.Assertf(amount > 0, "expense.Create: %s", InvalidExpenseAmountErr)
	assert.Assertf(len(label) > 0, "expense.Create: %s", InvalidExpenseLabelErr)

	stmt, err := db.Prepare(`
        insert into expenses (spender_id, label, amount, datetime, author_id, category_id)
        values (?, ?, ?, ?, ?, ?)
        returning id;
        `)
	assert.AssertErr(err)
	defer stmt.Close()

	categoryId := sql.NullInt32{Valid: false}
	if category != nil {
		categoryId = sql.NullInt32{Valid: true, Int32: int32(category.Id)}
	}

	if err = stmt.QueryRow(spender.Sub, label, amount, timestamp.Unix(), author.Sub, categoryId).Scan(&id); err != nil {
		return 0, fmt.Errorf("expense.Create: %w, %w", ExpenseCreateErr, err)
	}

	return id, nil
}

func (e Expense) Dto() dtos.Expense {
	return dtos.Expense{
		Id:         e.Id,
		SpenderId:  e.SpenderId,
		Label:      e.Label,
		Amount:     e.Amount,
		Timestamp:  dtos.UnixTime(e.Timestamp),
		AuthorId:   e.AuthorId,
		CategoryId: e.CategoryId,
	}
}
