package api


import (
	"encoding/json"
	"net/http"
)


func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func respondAccessError(w http.ResponseWriter, err error) {
	if err.Error() == "not found" {
		respondError(w, http.StatusNotFound, "character not found")
	} else {
		respondError(w, http.StatusForbidden, "you do not own this character")
	}
}
