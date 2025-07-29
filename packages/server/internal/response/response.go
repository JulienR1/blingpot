package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func Json(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("could not marshal data for request:", err)
		http.Error(w, "could not marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
