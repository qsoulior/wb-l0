package nats

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/stan.go"
	"github.com/qsoulior/wb-l0/internal/entity"
	"github.com/qsoulior/wb-l0/internal/service"
)

type handler struct {
	service service.Service
	logger  *slog.Logger
}

func NewHandler(s service.Service, logger *slog.Logger) *handler { return &handler{s, logger} }

func (h *handler) Serve(ctx context.Context) stan.MsgHandler {
	return func(msg *stan.Msg) {

		var order entity.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			h.logger.Error("json.Unmarshal", "err", err)
			return
		}

		_, err = h.service.Create(ctx, order)
		if err != nil {
			h.logger.Error("h.service.Create", "err", err)
			return
		}

		h.logger.Info("message received", "uid", order.OrderUID)
	}
}
