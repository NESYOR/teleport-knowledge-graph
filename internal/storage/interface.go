package storage

import (
	"context"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// SnapshotStore stores and retrieves snapshots.
type SnapshotStore interface {
	Save(ctx context.Context, snapshot model.Snapshot) error
	Load(ctx context.Context, id string) (model.Snapshot, error)
	LoadCurrent(ctx context.Context) (model.Snapshot, error)
}
