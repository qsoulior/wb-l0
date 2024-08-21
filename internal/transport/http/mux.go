package http

import (
	"log"
	"net/http"

	"github.com/qsoulior/wb-l0/internal/service"
)

func NewMux(s service.Service, logger *log.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", &handler{s, logger})
	return mux
}
