package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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
	client, auth, err := factory.New(timeoutCtx)
	if err != nil {
		logger.ErrorContext(timeoutCtx, "failed to initialize teleport client", "error", err)
		os.Exit(1)
	}

	orch := collectors.NewOrchestrator(logger,
		collectors.ClusterCollector{},
		collectors.UsersCollector{},
		collectors.RolesCollector{},
		collectors.NodesCollector{},
	)
	snapshot := orch.Collect(timeoutCtx, client, auth)

	store := storage.NewFSJSONStore(cfg.Storage.Path)
	if err := store.Save(timeoutCtx, snapshot); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.ErrorContext(timeoutCtx, "collection timed out before snapshot persisted", "error", err)
		} else {
			logger.ErrorContext(timeoutCtx, "failed to persist snapshot", "error", err)
		}
		os.Exit(1)
	}

	logger.InfoContext(timeoutCtx, "collection completed", "snapshot_id", snapshot.SnapshotID, "entities", len(snapshot.Entities), "relationships", len(snapshot.Relationships), "errors", len(snapshot.CollectorErrors))
	fmt.Println(snapshot.SnapshotID)
}
