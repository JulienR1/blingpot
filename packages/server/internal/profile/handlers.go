package profile

import (
	"log"
	"net/http"

	"github.com/julienr1/blingpot/internal/assert"
	"github.com/julienr1/blingpot/internal/database"
	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/response"
)

func HandleFindMe(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("profile").(Profile)
	response.Json(w, p.Dto())
}

func HandleFindAll(w http.ResponseWriter, r *http.Request) {
	db, err := database.Open()
	assert.AssertErr(err)
	defer db.Close()

	profiles, err := FindAll(db)
	if err != nil {
		log.Println("HandleFindAdd: could not fetch profiles:", err)
		http.Error(w, "could not fetch profiles", http.StatusInternalServerError)
		return
	}

	payload := make([]dtos.Profile, len(profiles))
	for i, p := range profiles {
		payload[i] = p.Dto()
	}

	response.Json(w, payload)
}
