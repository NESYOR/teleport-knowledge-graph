package collectors

import (
	"context"
	"errors"
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

type fakeClient struct{}

func (fakeClient) ClusterName(context.Context) (string, error) { return "root", nil }
func (fakeClient) CurrentUser(context.Context) (string, error) { return "alice", nil }
func (fakeClient) ListRoles(context.Context) ([]teleport.RoleRecord, error) {
	return []teleport.RoleRecord{{Name: "dev", Logins: []string{"ubuntu"}}}, nil
}
func (fakeClient) ListUsers(context.Context) ([]teleport.UserRecord, error) {
	return []teleport.UserRecord{{Name: "alice", Roles: []string{"dev"}}}, nil
}
func (fakeClient) ListNodes(context.Context) ([]teleport.NodeRecord, error) {
	return []teleport.NodeRecord{{Name: "n1", Namespace: "default", Labels: map[string]string{"env": "dev"}}}, nil
}

func TestOrchestratorCollect(t *testing.T) {
	o := NewOrchestrator(nil, ClusterCollector{}, UsersCollector{}, RolesCollector{}, NodesCollector{})
	s := o.Collect(context.Background(), fakeClient{}, teleport.AuthSource{Mode: "profile", Detail: "~/.tsh"})
	if s.Source.ClusterName != "root" {
		t.Fatalf("expected cluster root, got %s", s.Source.ClusterName)
	}
	if len(s.Entities) == 0 {
		t.Fatal("expected entities")
	}
	if len(s.Relationships) == 0 {
		t.Fatal("expected relationships")
	}
}

type failingClient struct{ fakeClient }

func (failingClient) ListUsers(context.Context) ([]teleport.UserRecord, error) {
	return nil, errors.New("permission denied")
}

func TestOrchestratorPartialError(t *testing.T) {
	o := NewOrchestrator(nil, UsersCollector{}, RolesCollector{})
	s := o.Collect(context.Background(), failingClient{}, teleport.AuthSource{Mode: "profile", Detail: "x"})
	if len(s.CollectorErrors) == 0 {
		t.Fatal("expected collector errors")
	}
	if len(s.Entities) == 0 {
		t.Fatal("expected role entities despite users collector failure")
	}
}
