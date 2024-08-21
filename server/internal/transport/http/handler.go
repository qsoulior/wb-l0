package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/qsoulior/wb-l0/internal/service"
)

type handler struct {
	service service.Service
	logger  *slog.Logger
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
		h.logger.Error("h.service.Get", "err", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

type page struct {
}

func (p *page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("../internal/transport/template/index.html")
	t.Execute(w, nil)
}
