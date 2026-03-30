package character

import (
	"log"
	"net/http"

	"github.com/cscercel/beyond-dnd/internal/json"
)


type handler struct {
	service Service	
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListCharacters(w http.ResponseWriter, r *http.Request) {
	err := h.service.ListCharacters(r.Context())

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	characters := struct {
		Characters []string `json:"characters"`
	}{}

	json.Write(w, http.StatusOK, characters)
}
