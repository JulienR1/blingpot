package expense

import (
	"errors"
	"fmt"
	"time"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/profile"
)

var InvalidExpenseAmountErr = errors.New("invalid expense amount")
var InvalidExpenseLabelErr = errors.New("invalid expense label")
var ExpenseCreateErr = errors.New("could not create expense")

type Expense struct {
	Id        int
	SpenderId string
	Label     string
	Amount    int
	Timestamp time.Time
	AuthorId  string
}

func FindById() (Expense, error) {
	return Expense{}, nil
}

func Create(db database.Querier, label string, amount int, timestamp time.Time, spender, author *profile.Profile) (id int, err error) {
	assert.Assertf(amount > 0, "expense.Create: %s", InvalidExpenseAmountErr)
	assert.Assertf(len(label) > 0, "expense.Create: %s", InvalidExpenseLabelErr)

	stmt, err := db.Prepare(`
        insert into expenses (spender_id, label, amount, datetime, author_id)
        values (?, ?, ?, ?, ?)
        returning id;
        `)
	assert.AssertErr(err)
	defer stmt.Close()

	if err = stmt.QueryRow(spender.Sub, label, amount, timestamp.Unix(), author.Sub).Scan(&id); err != nil {
		return 0, fmt.Errorf("expense.Create: %w, %w", ExpenseCreateErr, err)
	}

	return id, nil
}

func (e Expense) Dto() dtos.Expense {
	return dtos.Expense{}
}
