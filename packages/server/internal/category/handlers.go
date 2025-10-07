package category

import (
	"log"
	"net/http"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/response"
)

func HandleFindAll(w http.ResponseWriter, r *http.Request) {
	db, err := database.Open()
	assert.AssertErr(err)
	defer db.Close()

	categories, err := FindAll(db)
	if err != nil {
		log.Println("HandleFindAll: could not fetch categories:", err)
		http.Error(w, "could not fetch categories", http.StatusInternalServerError)
		return
	}

	payload := make([]dtos.Category, len(categories))
	for i, c := range categories {
		payload[i] = c.Dto()
	}

	response.Json(w, payload)
}
