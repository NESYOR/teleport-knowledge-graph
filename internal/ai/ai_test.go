package ai

import (
	"context"
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

func makeGraph(t *testing.T) *graph.Graph {
	t.Helper()
	g := graph.New()
	user := model.Entity{ID: model.DeterministicID("user", "alice"), Type: model.EntityUser, Name: "alice"}
	role := model.Entity{ID: model.DeterministicID("role", "prod-admin"), Type: model.EntityRole, Name: "prod-admin"}
	login := model.Entity{ID: model.DeterministicID("login", "root"), Type: model.EntityLogin, Name: "root"}
	node := model.Entity{ID: model.DeterministicID("node", "prod-node", "default"), Type: model.EntityNode, Name: "prod-node", Labels: map[string]string{"env": "prod"}}
	g.AddNode(user)
	g.AddNode(role)
	g.AddNode(login)
	g.AddNode(node)
	if err := g.AddEdge(model.Relationship{ID: model.DeterministicID("edge", "u2r", user.ID, role.ID), Type: model.RelUserHasRole, FromID: user.ID, ToID: role.ID, Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge(model.Relationship{ID: model.DeterministicID("edge", "r2n", role.ID, node.ID), Type: model.RelAccessPath, FromID: role.ID, ToID: node.ID, Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddEdge(model.Relationship{ID: model.DeterministicID("edge", "r2l", role.ID, login.ID), Type: model.RelRoleAllowsLogin, FromID: role.ID, ToID: login.ID, Confidence: 1}); err != nil {
		t.Fatal(err)
	}
	return g
}

func TestGraphContextBuilder(t *testing.T) {
	b := NewGraphContextBuilder(makeGraph(t))
	ctx, err := b.Build("alice", 10)
	if err != nil {
		t.Fatalf("build context: %v", err)
	}
	if len(ctx.Nodes) == 0 {
		t.Fatal("expected nodes")
	}
}

func TestGraphSummarizer(t *testing.T) {
	s := NewGraphSummarizer(makeGraph(t))
	userID := model.DeterministicID("user", "alice")
	sum, err := s.SummarizeUser(userID, 5)
	if err != nil {
		t.Fatalf("summarize user: %v", err)
	}
	if len(sum.Roles) == 0 {
		t.Fatal("expected roles")
	}
}

func TestLocalRuleReasoner(t *testing.T) {
	r := NewLocalRuleReasoner(makeGraph(t))
	out, err := r.WhyUserCanAccess(model.DeterministicID("user", "alice"), model.DeterministicID("node", "prod-node", "default"))
	if err != nil {
		t.Fatalf("reasoning error: %v", err)
	}
	if out.Conclusion != "access_path_found" {
		t.Fatalf("unexpected conclusion %s", out.Conclusion)
	}
}

func TestStaticAdapter(t *testing.T) {
	a := StaticAdapter{}
	resp, err := a.Generate(context.Background(), GenerateRequest{Prompt: "hello"})
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	if resp.Content == "" {
		t.Fatal("expected content")
	}
}

func TestPromptTemplate(t *testing.T) {
	tpl := AccessWhyTemplate{}
	text, err := tpl.Render(ContextBundle{Question: "why", SchemaVersion: "v1"})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if text == "" {
		t.Fatal("expected prompt text")
	}
}
