package summary

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/query"
	"github.com/julienr1/blingpot/internal/response"
)

func HandleExpenses(w http.ResponseWriter, r *http.Request) {
	var start, end dtos.UnixTime
	var err error

	errors.Join(err, query.UnixTime(r, "start", &start))
	errors.Join(err, query.UnixTime(r, "end", &end))
	errors.Join(err, query.Less(time.Time(start).Unix(), time.Time(end).Unix()))

	if err != nil {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	db, err := database.Open()
	assert.AssertErr(err)
	defer db.Close()

	total := GetTotal(db, time.Time(start), time.Time(end))
	subtotals := GetCategoryTotals(db, time.Time(start), time.Time(end))

	dto := dtos.ExpensesSummary{Total: total, Subtotals: subtotals}
	response.Json(w, dto)
}
