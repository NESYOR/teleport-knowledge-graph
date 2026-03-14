package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/teleport-cluster-digital-twin/internal/config"
	"github.com/example/teleport-cluster-digital-twin/internal/logging"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(os.Getenv("TWIN_CONFIG"))
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger := logging.New(cfg.Log.Level, cfg.Log.Format)
	logger.InfoContext(ctx, "twin bootstrap complete", "mode", cfg.Mode)
}
