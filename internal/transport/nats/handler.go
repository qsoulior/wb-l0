package nats

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/stan.go"
	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/internal/service"
)

type handler struct {
	service service.Service
	logger  *log.Logger
}

func NewHandler(s service.Service, logger *log.Logger) *handler { return &handler{s, logger} }

func (h *handler) Serve(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {
		var order entity.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			h.logger.Printf("json.Unmarshal: %s", err)
			return
		}

		_, err = h.service.Create(ctx, order)
		if err != nil {
			h.logger.Printf("h.service.Create: %s", err)
			return
		}
	}
}
