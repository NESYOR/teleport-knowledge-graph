package ai

import (
	"fmt"
	"sort"
	"strings"

	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/pkg/twinapi"
)

// GraphContextBuilder builds bounded AI context bundles from graph subgraphs.
type GraphContextBuilder struct {
	Graph *graph.Graph
}

// NewGraphContextBuilder creates a graph-backed context builder.
func NewGraphContextBuilder(g *graph.Graph) *GraphContextBuilder {
	return &GraphContextBuilder{Graph: g}
}

// Build creates a deterministic context bundle for a question.
func (b *GraphContextBuilder) Build(question string, maxNodes int) (ContextBundle, error) {
	if b == nil || b.Graph == nil {
		return ContextBundle{}, fmt.Errorf("context builder graph is nil")
	}
	if maxNodes <= 0 {
		maxNodes = 25
	}

	nodes, edges := b.Graph.Export()
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	sort.Slice(edges, func(i, j int) bool { return edges[i].ID < edges[j].ID })

	selected := make([]model.Entity, 0, min(maxNodes, len(nodes)))
	q := strings.ToLower(question)
	for _, n := range nodes {
		if len(selected) >= maxNodes {
			break
		}
		if q == "" || strings.Contains(strings.ToLower(n.Name), q) || strings.Contains(strings.ToLower(n.ID), q) {
			selected = append(selected, n)
		}
	}
	for _, n := range nodes {
		if len(selected) >= maxNodes {
			break
		}
		if !containsNode(selected, n.ID) {
			selected = append(selected, n)
		}
	}

	allowed := map[string]bool{}
	outNodes := make([]ContextNode, 0, len(selected))
	for _, n := range selected {
		allowed[n.ID] = true
		outNodes = append(outNodes, ContextNode{ID: n.ID, Type: string(n.Type), Name: n.Name})
	}

	outEdges := make([]ContextEdge, 0)
	for _, e := range edges {
		if allowed[e.FromID] && allowed[e.ToID] {
			outEdges = append(outEdges, ContextEdge{Type: string(e.Type), From: e.FromID, To: e.ToID})
		}
	}

	facts := []string{
		fmt.Sprintf("selected_nodes=%d", len(outNodes)),
		fmt.Sprintf("selected_edges=%d", len(outEdges)),
	}

	return ContextBundle{
		Question:      question,
		SchemaVersion: twinapi.Version,
		Nodes:         outNodes,
		Edges:         outEdges,
		Facts:         facts,
		Metadata:      map[string]any{"max_nodes": maxNodes},
	}, nil
}

func containsNode(nodes []model.Entity, id string) bool {
	for _, n := range nodes {
		if n.ID == id {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
