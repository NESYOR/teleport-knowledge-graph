package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// InMemoryCurrentStore keeps the latest snapshot in memory.
type InMemoryCurrentStore struct {
	mu       sync.RWMutex
	snapshot model.Snapshot
	hasValue bool
}

// NewInMemoryCurrentStore creates an in-memory store.
func NewInMemoryCurrentStore() *InMemoryCurrentStore {
	return &InMemoryCurrentStore{}
}

// Save stores snapshot as current.
func (s *InMemoryCurrentStore) Save(ctx context.Context, snapshot model.Snapshot) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshot = snapshot
	s.hasValue = true
	return nil
}

// Load returns snapshot only if it matches requested ID.
func (s *InMemoryCurrentStore) Load(ctx context.Context, id string) (model.Snapshot, error) {
	select {
	case <-ctx.Done():
		return model.Snapshot{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.hasValue || s.snapshot.SnapshotID != id {
		return model.Snapshot{}, errors.New("snapshot not found")
	}
	return s.snapshot, nil
}

// LoadCurrent returns current snapshot.
func (s *InMemoryCurrentStore) LoadCurrent(ctx context.Context) (model.Snapshot, error) {
	select {
	case <-ctx.Done():
		return model.Snapshot{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.hasValue {
		return model.Snapshot{}, errors.New("current snapshot not found")
	}
	return s.snapshot, nil
}
