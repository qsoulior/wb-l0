package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/qsoulior/wb-l0/internal/service"
)

type handler struct {
	service service.Service
	logger  *log.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		ErrorJSON(w, "order id is empty", http.StatusBadRequest)
		return
	}

	order, err := h.service.Get(r.Context(), id)
	if errors.Is(err, service.ErrNotExist) {
		ErrorJSON(w, "order does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		ErrorJSON(w, "internal server error", http.StatusInternalServerError)
		h.logger.Printf("h.service.Get: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
