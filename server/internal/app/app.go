package app

import (
	"context"
	"log/slog"
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

func Run(cfg *Config, logger *slog.Logger) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// database connection
	pg, err := postgres.New(ctx, cfg.Postgres.URI)
	if err != nil {
		logger.Error("failed to establish database connection", "err", err)
		return
	}
	defer pg.Close()
	logger.Info("database connection established", "uri", cfg.Postgres.URI)

	// repositories and service initialization
	db := repo.NewPG(pg)
	cache := repo.NewCache(ctx)
	service := service.New(db, cache)
	err = service.Init(ctx)
	if err != nil {
		logger.Error("failed to init service", "err", err)
		return
	}

	// http server
	mux := http.NewMux(service, logger)
	server := httpserver.New(mux, cfg.HTTP.Host, cfg.HTTP.Port)
	server.Start(ctx)
	logger.Info("http server started", "host", cfg.HTTP.Host, "port", cfg.HTTP.Port)

	// nats streaming
	scErrCh := make(chan error, 1)
	sc, err := stan.Connect(cfg.NATS.ClusterID, cfg.NATS.ClientID,
		stan.Pings(30, 20),
		stan.NatsURL(cfg.NATS.URL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, err error) {
			scErrCh <- err
			close(scErrCh)
		}))
	if err != nil {
		logger.Error("failed to establish stan connection", "err", err)
		return
	}
	logger.Info("stan connection established", "url", cfg.NATS.URL)
	sh := nats.NewHandler(service, logger).Serve(ctx)
	sub, err := sc.Subscribe(cfg.NATS.Subscription, sh, stan.DurableName("app-durable"))
	if err != nil {
		logger.Error("failed to subscribe", "err", err)
		return
	}
	logger.Info("stan subscription started")

	// graceful shutdown
	select {
	case <-ctx.Done():
		logger.Info("notify signal accepted")
	case err := <-server.Err():
		logger.Error("http server returned error", "err", err)
	case err := <-scErrCh:
		logger.Error("stan connection handler returns error", "err", err)
	}

	// http server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Stop(ctx)
	if err != nil {
		logger.Error("failed to stop http server", "err", err)
	}

	// stan shutdown
	err = sub.Unsubscribe()
	if err != nil {
		logger.Error("failed to unsubscribe", "err", err)
	}

	err = sc.Close()
	if err != nil {
		logger.Error("failed to close stan connection", "err", err)
	}
}
