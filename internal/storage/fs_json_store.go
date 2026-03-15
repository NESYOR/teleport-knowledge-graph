package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// FSJSONStore stores snapshots as JSON files in a directory.
type FSJSONStore struct {
	basePath string
}

// NewFSJSONStore creates a filesystem-backed snapshot store.
func NewFSJSONStore(basePath string) *FSJSONStore {
	return &FSJSONStore{basePath: basePath}
}

// Save writes a snapshot to disk and updates current.json.
func (s *FSJSONStore) Save(ctx context.Context, snapshot model.Snapshot) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if err := os.MkdirAll(s.basePath, 0o755); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	b, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	file := filepath.Join(s.basePath, snapshot.SnapshotID+".json")
	if err := os.WriteFile(file, b, 0o644); err != nil {
		return fmt.Errorf("write snapshot file: %w", err)
	}

	current := filepath.Join(s.basePath, "current.json")
	if err := os.WriteFile(current, b, 0o644); err != nil {
		return fmt.Errorf("write current snapshot: %w", err)
	}
	return nil
}

// Load reads a snapshot by ID.
func (s *FSJSONStore) Load(ctx context.Context, id string) (model.Snapshot, error) {
	select {
	case <-ctx.Done():
		return model.Snapshot{}, ctx.Err()
	default:
	}
	file := filepath.Join(s.basePath, id+".json")
	return s.read(file)
}

// LoadCurrent reads current snapshot pointer.
func (s *FSJSONStore) LoadCurrent(ctx context.Context) (model.Snapshot, error) {
	select {
	case <-ctx.Done():
		return model.Snapshot{}, ctx.Err()
	default:
	}
	file := filepath.Join(s.basePath, "current.json")
	return s.read(file)
}

func (s *FSJSONStore) read(file string) (model.Snapshot, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return model.Snapshot{}, fmt.Errorf("read snapshot: %w", err)
	}
	var out model.Snapshot
	if err := json.Unmarshal(b, &out); err != nil {
		return model.Snapshot{}, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return out, nil
}
