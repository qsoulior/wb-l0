package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/qsoulior/wb-l0/internal/repo"
	"github.com/qsoulior/wb-l0/internal/service"
	"github.com/qsoulior/wb-l0/internal/transport/http"
	"github.com/qsoulior/wb-l0/internal/transport/nats"
	"github.com/qsoulior/wb-l0/pkg/httpserver"
	"github.com/qsoulior/wb-l0/pkg/postgres"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// database connection
	pg, err := postgres.New(ctx, "TODO")
	if err != nil {
		logger.Error("failed to establish database connection", "err", err)
		return
	}
	defer pg.Close()

	// repositories and service initialization
	db := repo.NewPG(pg)
	cache := repo.NewCache(ctx)
	service := service.New(db, cache)

	// http server
	mux := http.NewMux(service, logger)
	server := httpserver.New(mux, "TODO", "TODO")
	server.Start(ctx)
	logger.Info("http server started", "host", "TODO", "port", "TODO")

	// nats streaming
	sc, err := stan.Connect("TODO", "TODO")
	if err != nil {
		logger.Error("failed to establish stan connection", "err", err)
		return
	}
	sh := nats.NewHandler(service, logger).Serve(ctx)
	sub, err := sc.Subscribe("", sh)
	if err != nil {
		logger.Error("failed to subscribe", "err", err)
		return
	}
	logger.Info("stan subscription started")

	// graceful shutdown
	select {
	case <-ctx.Done():
		logger.Info("notify signal accepted")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = server.Stop(ctx)
		if err != nil {
			logger.Error("failed to stop http server", "err", err)
		}

		err = sub.Unsubscribe()
		if err != nil {
			logger.Error("failed to unsubscribe", "err", err)
		}

		err = sc.Close()
		if err != nil {
			logger.Error("failed to close stan connection", "err", err)
		}
	case err := <-server.Err():
		logger.Error("http server returned error", "err", err)
	}
}
