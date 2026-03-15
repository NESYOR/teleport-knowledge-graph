package graph

import "github.com/example/teleport-cluster-digital-twin/internal/model"

// Node gets a node by ID.
func (g *Graph) Node(id string) (model.Entity, bool) {
	n, ok := g.nodesByID[id]
	return n, ok
}

// NodesByType returns all nodes for a specific entity type.
func (g *Graph) NodesByType(t model.EntityType) []model.Entity {
	ids := g.nodeTypeIdx[t]
	out := make([]model.Entity, 0, len(ids))
	for _, id := range ids {
		out = append(out, g.nodesByID[id])
	}
	return out
}

// OutgoingEdges returns outgoing edges for a node.
func (g *Graph) OutgoingEdges(id string) []model.Relationship {
	eids := g.outgoing[id]
	out := make([]model.Relationship, 0, len(eids))
	for _, eid := range eids {
		out = append(out, g.edgesByID[eid])
	}
	return out
}

// IncomingEdges returns incoming edges for a node.
func (g *Graph) IncomingEdges(id string) []model.Relationship {
	eids := g.incoming[id]
	out := make([]model.Relationship, 0, len(eids))
	for _, eid := range eids {
		out = append(out, g.edgesByID[eid])
	}
	return out
}
