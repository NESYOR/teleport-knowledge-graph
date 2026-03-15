package collectors

import (
	"context"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// ClusterCollector captures cluster identity as an entity.
type ClusterCollector struct{}

// Name returns collector name.
func (ClusterCollector) Name() string { return "cluster" }

// Collect returns the cluster entity.
func (ClusterCollector) Collect(ctx context.Context, client teleport.ResourceClient) Result {
	cluster, err := client.ClusterName(ctx)
	if err != nil {
		return Result{Errors: []model.CollectorError{{Collector: "cluster", Kind: teleport.ErrorKind(err), Message: err.Error()}}}
	}

	e := model.Entity{
		ID:      model.DeterministicID("cluster", cluster),
		Type:    model.EntityCluster,
		Name:    cluster,
		Cluster: cluster,
	}
	return Result{Entities: []model.Entity{e}, Counts: map[string]int{"clusters": 1}}
}
