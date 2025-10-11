package expense

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/category"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/profile"
	"github.com/julienr1/blingpot/internal/query"
	"github.com/julienr1/blingpot/internal/response"
)

type FindParams struct {
	Start dtos.UnixTime
	End   dtos.UnixTime
}

type CreateExpenseBody struct {
	Label      string `json:"label" validate:"required,min=1"`
	Amount     int    `json:"amount" validate:"required,number,gt=0"`
	SpenderId  string `json:"spenderId" validate:"required,min=1,alphanum"`
	Timestamp  int64  `json:"timestamp" validate:"required,number"`
	CategoryId *int   `json:"categoryId"`
}

type CreateResponseBody struct {
	Id int `json:"id"`
}

func HandleFind(w http.ResponseWriter, r *http.Request) {
	var start, end dtos.UnixTime
	var err error

	errors.Join(err, query.UnixTime(r, "start", &start))
	errors.Join(err, query.UnixTime(r, "end", &end))

	if err != nil {
		http.Error(w, "invalid request parameters", http.StatusBadRequest)
		return
	}

	db, err := database.Open()
	assert.AssertErr(err)
	defer db.Close()

	expenses, err := Find(db, time.Time(start), time.Time(end))
	if err != nil {
		log.Printf("could not find expense: %s\r\n", err.Error())
		http.Error(w, "could not find expense", http.StatusInternalServerError)
		return
	}

	var dtos = make([]dtos.Expense, len(expenses))
	for i, expense := range expenses {
		dtos[i] = expense.Dto()
	}

	response.Json(w, dtos)
}

func HandleCreate(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("profile").(profile.Profile)

	var body CreateExpenseBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("could not read body: %s", err), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %s", err), http.StatusBadRequest)
		return
	}

	db, err := database.Open()
	assert.AssertErr(err)
	defer db.Close()

	spender, err := profile.FindBySub(db, body.SpenderId)
	if err != nil {
		http.Error(w, "specified spender profile does not exist", http.StatusBadRequest)
		return
	}

	var c *category.Category = nil
	if body.CategoryId != nil {
		if c, err = category.FindById(db, *body.CategoryId); err != nil {
			http.Error(w, "specified category does not exist", http.StatusBadRequest)
		}
	}

	timestamp := time.Unix(body.Timestamp, 0)
	id, err := Create(db, body.Label, body.Amount, timestamp, spender, &p, c)
	if err != nil {
		log.Printf("could not create expense: %s\r\n", err.Error())
		http.Error(w, "could not create expense", http.StatusInternalServerError)
		return
	}

	response.Json(w, CreateResponseBody{Id: id})
}
