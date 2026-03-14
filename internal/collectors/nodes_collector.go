package collectors

import (
	"context"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// NodesCollector collects node resources.
type NodesCollector struct{}

// Name returns collector name.
func (NodesCollector) Name() string { return "nodes" }

// Collect gathers node entities and normalized labels.
func (NodesCollector) Collect(ctx context.Context, client teleport.ResourceClient) Result {
	nodes, err := client.ListNodes(ctx)
	if err != nil {
		return Result{Errors: []model.CollectorError{{Collector: "nodes", Kind: teleport.ErrorKind(err), Message: err.Error()}}}
	}
	entities := make([]model.Entity, 0, len(nodes))
	labelCount := 0
	for _, n := range nodes {
		entities = append(entities, model.Entity{
			ID:        model.DeterministicID("node", n.Name, n.Namespace),
			Type:      model.EntityNode,
			Name:      n.Name,
			Namespace: n.Namespace,
			Labels:    n.Labels,
		})
		labelCount += len(n.Labels)
	}
	return Result{Entities: entities, Counts: map[string]int{"nodes": len(nodes), "labels": labelCount}}
}
