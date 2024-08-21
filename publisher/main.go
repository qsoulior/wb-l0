package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

type Config struct {
	URL       string `json:"url"`
	ClusterID string `json:"cluster_id"`
	ClientID  string `json:"client_id"`
	Subject   string `json:"subject"`
	Interval  int    `json:"interval"`
	Count     int    `json:"count"`
}

var (
	DefaultInterval = 10 * time.Second
	DefaultFile     = "model.json"
)

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	d := json.NewDecoder(f)
	err = d.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	var path string
	flag.StringVar(&path, "c", "", "config file path")
	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := NewConfig(path)
	if err != nil {
		logger.Error("failed to load config", "err", err)
		return
	}

	f, err := os.Open(DefaultFile)
	if err != nil {
		logger.Error("failed to load model file", "err", err)
		return
	}

	order := new(Order)
	d := json.NewDecoder(f)
	err = d.Decode(order)
	if err != nil {
		logger.Error("failed to decode model", "err", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsURL(cfg.URL))
	if err != nil {
		logger.Error("failed to establish stan connection", "err", err)
		return
	}
	logger.Info("stan connection established", "url", cfg.URL)

	interval := DefaultInterval
	if cfg.Interval > 0 {
		interval = time.Duration(cfg.Interval) * time.Second
	}

	timer := time.NewTimer(0)
	for count := 0; count < cfg.Count || cfg.Count <= 0; count++ {
		select {
		case <-timer.C:
			order.OrderUID = uuid.New().String()
			data, err := json.Marshal(order)
			if err != nil {
				logger.Error("failed to encode message", "err", err)
				continue
			}

			err = sc.Publish(cfg.Subject, data)
			if err != nil {
				logger.Error("failed to publish message", "err", err)
				continue
			}
			logger.Info("message sended", "uid", order.OrderUID)
			timer.Reset(interval)
		case <-ctx.Done():
			logger.Info("notify signal accepted")
			timer.Stop()
			return
		}
	}
}
