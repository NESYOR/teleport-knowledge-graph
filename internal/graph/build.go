package graph

import "github.com/example/teleport-cluster-digital-twin/internal/model"

// FromSnapshot builds a graph instance from snapshot entities and relationships.
func FromSnapshot(s model.Snapshot) *Graph {
	g := New()
	for _, n := range s.Entities {
		g.AddNode(n)
	}
	for _, e := range s.Relationships {
		_ = g.AddEdge(e)
	}
	return g
}
