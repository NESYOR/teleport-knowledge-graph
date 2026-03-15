package graph

import (
	"fmt"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// Graph is a typed in-memory graph with forward and reverse adjacency.
type Graph struct {
	nodesByID   map[string]model.Entity
	edgesByID   map[string]model.Relationship
	outgoing    map[string][]string
	incoming    map[string][]string
	nodeTypeIdx map[model.EntityType][]string
}

// New returns an initialized graph.
func New() *Graph {
	return &Graph{
		nodesByID:   map[string]model.Entity{},
		edgesByID:   map[string]model.Relationship{},
		outgoing:    map[string][]string{},
		incoming:    map[string][]string{},
		nodeTypeIdx: map[model.EntityType][]string{},
	}
}

// AddNode inserts or replaces a node.
func (g *Graph) AddNode(n model.Entity) {
	if _, ok := g.nodesByID[n.ID]; !ok {
		g.nodeTypeIdx[n.Type] = append(g.nodeTypeIdx[n.Type], n.ID)
	}
	g.nodesByID[n.ID] = n
}

// AddEdge inserts an edge and updates adjacency lists.
func (g *Graph) AddEdge(e model.Relationship) error {
	if _, ok := g.nodesByID[e.FromID]; !ok {
		return fmt.Errorf("edge from node %q does not exist", e.FromID)
	}
	if _, ok := g.nodesByID[e.ToID]; !ok {
		return fmt.Errorf("edge to node %q does not exist", e.ToID)
	}
	if _, ok := g.edgesByID[e.ID]; ok {
		return nil
	}
	g.edgesByID[e.ID] = e
	g.outgoing[e.FromID] = append(g.outgoing[e.FromID], e.ID)
	g.incoming[e.ToID] = append(g.incoming[e.ToID], e.ID)
	return nil
}
