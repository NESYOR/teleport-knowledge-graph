package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/teleport-cluster-digital-twin/internal/api"
	"github.com/example/teleport-cluster-digital-twin/internal/cli"
	"github.com/example/teleport-cluster-digital-twin/internal/collectors"
	"github.com/example/teleport-cluster-digital-twin/internal/config"
	"github.com/example/teleport-cluster-digital-twin/internal/logging"
	"github.com/example/teleport-cluster-digital-twin/internal/storage"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
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
	timeoutCtx, cancel := context.WithTimeout(ctx, cfg.Collection.Timeout)
	defer cancel()

	factory := teleport.NewFactory(cfg.Teleport, logger)
	orch := collectors.NewOrchestrator(logger,
		collectors.ClusterCollector{},
		collectors.UsersCollector{},
		collectors.RolesCollector{},
		collectors.NodesCollector{},
	)
	store := storage.NewFSJSONStore(cfg.Storage.Path)

	if len(os.Args) > 1 && os.Args[1] == "serve" {
		srv := api.NewServer(factory, orch, store)
		httpServer := &http.Server{Addr: cfg.API.ListenAddr, Handler: srv.Routes()}
		logger.InfoContext(timeoutCtx, "starting API server", "addr", cfg.API.ListenAddr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(timeoutCtx, "api server failed", "error", err)
			os.Exit(1)
		}
		return
	}

	r := &cli.Runner{Factory: factory, Orchestrator: orch, Store: store, Out: os.Stdout}
	if err := r.Run(timeoutCtx, os.Args[1:]); err != nil {
		logger.ErrorContext(timeoutCtx, "command failed", "error", err)
		os.Exit(1)
	}
}
