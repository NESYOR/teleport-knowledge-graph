package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/collectors"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/storage"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

type fakeFactory struct{}

func (fakeFactory) New(context.Context) (teleport.ResourceClient, teleport.AuthSource, error) {
	return &teleport.LocalStaticClient{}, teleport.AuthSource{Mode: "profile", Detail: "test"}, nil
}

func newTestServer(t *testing.T) (*Server, *storage.FSJSONStore) {
	t.Helper()
	dir := t.TempDir()
	store := storage.NewFSJSONStore(filepath.Join(dir, "snapshots"))
	orch := collectors.NewOrchestrator(nil, collectors.ClusterCollector{}, collectors.UsersCollector{}, collectors.RolesCollector{}, collectors.NodesCollector{})
	s := NewServer(fakeFactory{}, orch, store)
	return s, store
}

func TestCollectAndSummary(t *testing.T) {
	s, _ := newTestServer(t)
	r := httptest.NewRequest(http.MethodPost, "/collect", nil)
	w := httptest.NewRecorder()
	s.Routes().ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	r = httptest.NewRequest(http.MethodGet, "/summary", nil)
	w = httptest.NewRecorder()
	s.Routes().ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestDiffEndpoint(t *testing.T) {
	s, store := newTestServer(t)
	old := model.Snapshot{SnapshotID: "old", Entities: []model.Entity{{ID: "user:1", Type: model.EntityUser, Name: "alice"}}}
	newS := model.Snapshot{SnapshotID: "new", Entities: []model.Entity{{ID: "user:1", Type: model.EntityUser, Name: "alice2"}}}
	if err := store.Save(context.Background(), old); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(context.Background(), newS); err != nil {
		t.Fatal(err)
	}

	payload, _ := json.Marshal(map[string]string{"old_id": "old", "new_id": "new"})
	r := httptest.NewRequest(http.MethodPost, "/diff", bytes.NewReader(payload))
	w := httptest.NewRecorder()
	s.Routes().ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}
