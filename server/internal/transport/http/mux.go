package http

import (
	"log/slog"
	"net/http"

	"github.com/qsoulior/wb-l0/internal/service"
)

func NewMux(s service.Service, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", &handler{s, logger})
	return mux
}
