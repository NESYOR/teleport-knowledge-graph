package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/collectors"
	"github.com/example/teleport-cluster-digital-twin/internal/storage"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

type fakeFactory struct{}

func (fakeFactory) New(context.Context) (teleport.ResourceClient, teleport.AuthSource, error) {
	return &teleport.LocalStaticClient{}, teleport.AuthSource{Mode: "profile", Detail: "test"}, nil
}

func TestRunnerCollectAndShowSummary(t *testing.T) {
	dir := t.TempDir()
	store := storage.NewFSJSONStore(filepath.Join(dir, "snapshots"))
	orch := collectors.NewOrchestrator(nil, collectors.ClusterCollector{}, collectors.UsersCollector{}, collectors.RolesCollector{}, collectors.NodesCollector{})

	buf := &bytes.Buffer{}
	r := &Runner{Factory: fakeFactory{}, Orchestrator: orch, Store: store, Out: buf}
	if err := r.Run(context.Background(), []string{"collect"}); err != nil {
		t.Fatalf("collect failed: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Fatal("expected snapshot id output")
	}

	buf.Reset()
	if err := r.Run(context.Background(), []string{"show", "summary"}); err != nil {
		t.Fatalf("show summary failed: %v", err)
	}
	if !strings.Contains(buf.String(), "entity_counts") {
		t.Fatalf("expected summary json output, got: %s", buf.String())
	}
}
