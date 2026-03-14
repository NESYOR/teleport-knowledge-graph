package graph

import "github.com/example/teleport-cluster-digital-twin/internal/model"

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
	return nodes, edges
}
