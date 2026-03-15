package graph

import (
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

func TestExportDeterministicOrder(t *testing.T) {
	g := New()
	a := model.Entity{ID: "b", Type: model.EntityUser, Name: "b"}
	b := model.Entity{ID: "a", Type: model.EntityUser, Name: "a"}
	g.AddNode(a)
	g.AddNode(b)
	if err := g.AddEdge(model.Relationship{ID: "e2", Type: model.RelUserHasRole, FromID: "a", ToID: "b", Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge(model.Relationship{ID: "e1", Type: model.RelUserHasRole, FromID: "b", ToID: "a", Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	nodes, edges := g.Export()
	if nodes[0].ID != "a" || nodes[1].ID != "b" {
		t.Fatalf("nodes not sorted: %#v", nodes)
	}
	if edges[0].ID != "e1" || edges[1].ID != "e2" {
		t.Fatalf("edges not sorted: %#v", edges)
	}
}
