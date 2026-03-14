package analysis

import "github.com/example/teleport-cluster-digital-twin/internal/graph"

// Service provides analysis operations over a graph.
type Service struct {
	Graph *graph.Graph
}

// New creates an analysis service.
func New(g *graph.Graph) *Service {
	return &Service{Graph: g}
}
