package graph

import (
	"sort"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// Export converts graph contents into snapshot entities and relationships.
func (g *Graph) Export() ([]model.Entity, []model.Relationship) {
	nodes := make([]model.Entity, 0, len(g.nodesByID))
	edges := make([]model.Relationship, 0, len(g.edgesByID))
	for _, n := range g.nodesByID {
		nodes = append(nodes, n)
	}
	for _, e := range g.edgesByID {
		edges = append(edges, e)
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	sort.Slice(edges, func(i, j int) bool { return edges[i].ID < edges[j].ID })
	return nodes, edges
}
