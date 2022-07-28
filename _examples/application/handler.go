package application

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	api CustomerAPI
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	id := parts[len(parts)-1]

	customer, err := h.api.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(customer)
	if err != nil {
		log.Println(err)
	}
}
