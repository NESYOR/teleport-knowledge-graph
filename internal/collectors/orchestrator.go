package collectors

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// Orchestrator runs collectors, captures partial failures, and assembles snapshots.
type Orchestrator struct {
	collectors []Collector
	logger     *slog.Logger
}

// NewOrchestrator creates a collection orchestrator.
func NewOrchestrator(logger *slog.Logger, collectors ...Collector) *Orchestrator {
	if logger == nil {
		logger = slog.Default()
	}
	return &Orchestrator{collectors: collectors, logger: logger}
}

// Collect executes configured collectors and returns a normalized snapshot.
func (o *Orchestrator) Collect(ctx context.Context, client teleport.ResourceClient, auth teleport.AuthSource) model.Snapshot {
	start := time.Now().UTC()
	snapshot := model.Snapshot{
		SnapshotID:      model.DeterministicID("snapshot", start.Format(time.RFC3339), auth.Mode, auth.Detail),
		CollectedAt:     start,
		CollectorStatus: make([]model.CollectorStatus, 0, len(o.collectors)),
		Entities:        []model.Entity{},
		Relationships:   []model.Relationship{},
		Metadata:        map[string]string{"auth_mode": auth.Mode, "auth_detail": auth.Detail},
	}

	if clusterName, err := client.ClusterName(ctx); err == nil {
		snapshot.Source.ClusterName = clusterName
	} else {
		snapshot.CollectorErrors = append(snapshot.CollectorErrors, model.CollectorError{Collector: "cluster_name", Kind: teleport.ErrorKind(err), Message: err.Error()})
	}
	if caller, err := client.CurrentUser(ctx); err == nil {
		snapshot.Source.Caller = caller
	} else {
		snapshot.CollectorErrors = append(snapshot.CollectorErrors, model.CollectorError{Collector: "current_user", Kind: teleport.ErrorKind(err), Message: err.Error()})
	}
	snapshot.Source.AuthSource = auth.Mode

	for _, c := range o.collectors {
		cStart := time.Now().UTC()
		result := c.Collect(ctx, client)
		cEnd := time.Now().UTC()

		snapshot.Entities = append(snapshot.Entities, result.Entities...)
		snapshot.Relationships = append(snapshot.Relationships, result.Relationships...)
		snapshot.CollectorErrors = append(snapshot.CollectorErrors, result.Errors...)
		snapshot.CollectorStatus = append(snapshot.CollectorStatus, model.CollectorStatus{
			Name:       c.Name(),
			StartedAt:  cStart,
			FinishedAt: cEnd,
			DurationMS: cEnd.Sub(cStart).Milliseconds(),
			Counts:     result.Counts,
		})
	}

	sort.Slice(snapshot.Entities, func(i, j int) bool { return snapshot.Entities[i].ID < snapshot.Entities[j].ID })
	sort.Slice(snapshot.Relationships, func(i, j int) bool { return snapshot.Relationships[i].ID < snapshot.Relationships[j].ID })
	sort.Slice(snapshot.CollectorErrors, func(i, j int) bool {
		if snapshot.CollectorErrors[i].Collector == snapshot.CollectorErrors[j].Collector {
			return snapshot.CollectorErrors[i].Message < snapshot.CollectorErrors[j].Message
		}
		return snapshot.CollectorErrors[i].Collector < snapshot.CollectorErrors[j].Collector
	})

	o.logger.Info("collection completed", "collectors", len(o.collectors), "entities", len(snapshot.Entities), "relationships", len(snapshot.Relationships), "errors", len(snapshot.CollectorErrors))
	return snapshot
}
