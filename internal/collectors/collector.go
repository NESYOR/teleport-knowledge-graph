package collectors

import (
	"context"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// Collector collects a resource slice and normalizes entities/relationships.
type Collector interface {
	Name() string
	Collect(ctx context.Context, client teleport.ResourceClient) Result
}

// Result is collector output with partial failure support.
type Result struct {
	Entities      []model.Entity
	Relationships []model.Relationship
	Errors        []model.CollectorError
	Counts        map[string]int
}
