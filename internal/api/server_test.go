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

func TestMethodGuards(t *testing.T) {
	s, _ := newTestServer(t)
	cases := []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/healthz"},
		{method: http.MethodGet, path: "/collect"},
		{method: http.MethodPost, path: "/summary"},
		{method: http.MethodPost, path: "/analysis/risks"},
		{method: http.MethodPost, path: "/roles/prod-admin"},
	}
	for _, tc := range cases {
		r := httptest.NewRequest(tc.method, tc.path, nil)
		w := httptest.NewRecorder()
		s.Routes().ServeHTTP(w, r)
		if w.Code != http.StatusMethodNotAllowed {
			t.Fatalf("expected 405 for %s %s got %d", tc.method, tc.path, w.Code)
		}
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

func TestDiffValidation(t *testing.T) {
	s, _ := newTestServer(t)
	payload, _ := json.Marshal(map[string]string{"old_id": "", "new_id": "new"})
	r := httptest.NewRequest(http.MethodPost, "/diff", bytes.NewReader(payload))
	w := httptest.NewRecorder()
	s.Routes().ServeHTTP(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}
