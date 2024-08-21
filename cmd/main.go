package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/qsoulior/wb-l0/internal/app"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "config file path")
	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := app.NewConfig(path)
	if err != nil {
		logger.Error("failed to load config", "err", err)
		return
	}

	app.Run(cfg, logger)
}
