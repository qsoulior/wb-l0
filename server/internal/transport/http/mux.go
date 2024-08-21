package http

import (
	"log/slog"
	"net/http"

	"github.com/qsoulior/wb-l0/internal/service"
)

func NewMux(s service.Service, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", &handler{s, logger})
	mux.Handle("GET /page", &page{logger})
	return mux
}
