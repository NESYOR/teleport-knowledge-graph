package analysis

import (
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

func testGraph(t *testing.T) *graph.Graph {
	t.Helper()
	g := graph.New()
	user := model.Entity{ID: model.DeterministicID("user", "alice"), Type: model.EntityUser, Name: "alice"}
	role := model.Entity{ID: model.DeterministicID("role", "prod-admin"), Type: model.EntityRole, Name: "prod-admin"}
	login := model.Entity{ID: model.DeterministicID("login", "root"), Type: model.EntityLogin, Name: "root"}
	node := model.Entity{ID: model.DeterministicID("node", "node-1", "default"), Type: model.EntityNode, Name: "node-1", Labels: map[string]string{"env": "*"}}
	for _, n := range []model.Entity{user, role, login, node} {
		g.AddNode(n)
	}
	mustEdge := func(e model.Relationship) {
		t.Helper()
		if err := g.AddEdge(e); err != nil {
			t.Fatalf("add edge: %v", err)
		}
	}
	mustEdge(model.Relationship{ID: model.DeterministicID("edge", "u2r", user.ID, role.ID), Type: model.RelUserHasRole, FromID: user.ID, ToID: role.ID, Confidence: 1})
	mustEdge(model.Relationship{ID: model.DeterministicID("edge", "r2n", role.ID, node.ID), Type: model.RelAccessPath, FromID: role.ID, ToID: node.ID, Confidence: 1})
	mustEdge(model.Relationship{ID: model.DeterministicID("edge", "r2l", role.ID, login.ID), Type: model.RelRoleAllowsLogin, FromID: role.ID, ToID: login.ID, Confidence: 1})
	return g
}

func TestExplainUserAccess(t *testing.T) {
	svc := New(testGraph(t))
	got, err := svc.ExplainUserAccess("alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got.Resources))
	}
}

func TestExplainResourceExposure(t *testing.T) {
	svc := New(testGraph(t))
	resourceID := model.DeterministicID("node", "node-1", "default")
	got, err := svc.ExplainResourceExposure(resourceID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Users) != 1 || got.Users[0] != "alice" {
		t.Fatalf("unexpected users: %#v", got.Users)
	}
}

func TestAnalyzeRisks(t *testing.T) {
	svc := New(testGraph(t))
	findings := svc.AnalyzeRisks()
	if len(findings) < 2 {
		t.Fatalf("expected at least 2 findings, got %d", len(findings))
	}
}

func TestTopology(t *testing.T) {
	svc := New(testGraph(t))
	top := svc.Topology()
	if top.EntityCounts[model.EntityUser] != 1 {
		t.Fatalf("expected one user, got %d", top.EntityCounts[model.EntityUser])
	}
	if top.TotalEdges != 3 {
		t.Fatalf("expected three edges, got %d", top.TotalEdges)
	}
}
