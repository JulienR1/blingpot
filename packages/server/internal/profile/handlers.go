package profile

import (
	"net/http"

	"github.com/julienr1/blingpot/internal/dtos"
	"github.com/julienr1/blingpot/internal/response"
)

func HandleFindMe(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("profile").(Profile)

	response.Json(w,
		dtos.Profile{
			FirstName: p.FirstName,
			LastName:  p.LastName,
			Email:     p.Email,
			Picture:   dtos.NullString(p.Picture),
		},
	)
}
