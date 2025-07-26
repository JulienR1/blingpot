package profile

import (
	"fmt"
	"net/http"
)

func HandleFindMe(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Context().Value("profile").(Profile))

	w.WriteHeader(http.StatusOK)
}
