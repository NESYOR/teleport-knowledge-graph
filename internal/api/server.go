package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/example/teleport-cluster-digital-twin/internal/analysis"
	"github.com/example/teleport-cluster-digital-twin/internal/collectors"
	"github.com/example/teleport-cluster-digital-twin/internal/diff"
	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/storage"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
	"github.com/example/teleport-cluster-digital-twin/pkg/twinapi"
)

// Server exposes HTTP handlers for snapshot, analysis and diff workflows.
type Server struct {
	Factory      teleport.ClientFactory
	Orchestrator *collectors.Orchestrator
	Store        storage.SnapshotStore
}

// NewServer creates a new API server.
func NewServer(factory teleport.ClientFactory, orch *collectors.Orchestrator, store storage.SnapshotStore) *Server {
	return &Server{Factory: factory, Orchestrator: orch, Store: store}
}

// Routes creates an HTTP mux with all API endpoints.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/collect", s.collect)
	mux.HandleFunc("/snapshot/current", s.current)
	mux.HandleFunc("/summary", s.summary)
	mux.HandleFunc("/users/", s.userAccess)
	mux.HandleFunc("/resources/", s.resourceExposure)
	mux.HandleFunc("/roles/", s.roleByName)
	mux.HandleFunc("/analysis/risks", s.risks)
	mux.HandleFunc("/diff", s.diff)
	return mux
}

func (s *Server) write(w http.ResponseWriter, code int, data interface{}, errs ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(Envelope{RequestID: time.Now().UTC().Format(time.RFC3339Nano), Version: twinapi.Version, Data: data, Errors: errs})
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.write(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}
	s.write(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) collect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.write(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}
	client, auth, err := s.Factory.New(r.Context())
	if err != nil {
		s.write(w, http.StatusBadGateway, nil, err.Error())
		return
	}
	snapshot := s.Orchestrator.Collect(r.Context(), client, auth)
	if err := s.Store.Save(r.Context(), snapshot); err != nil {
		s.write(w, http.StatusInternalServerError, nil, err.Error())
		return
	}
	s.write(w, http.StatusOK, snapshot)
}

func (s *Server) current(w http.ResponseWriter, r *http.Request) {
	snap, err := s.Store.LoadCurrent(r.Context())
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	s.write(w, http.StatusOK, snap)
}

func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	snap, err := s.Store.LoadCurrent(r.Context())
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	g := graph.FromSnapshot(snap)
	a := analysis.New(g)
	s.write(w, http.StatusOK, a.Topology())
}

func (s *Server) userAccess(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/access") {
		s.write(w, http.StatusNotFound, nil, "not found")
		return
	}
	name := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/users/"), "/access")
	snap, err := s.Store.LoadCurrent(r.Context())
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	out, err := analysis.New(graph.FromSnapshot(snap)).ExplainUserAccess(name)
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	s.write(w, http.StatusOK, out)
}

func (s *Server) resourceExposure(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/exposure") {
		s.write(w, http.StatusNotFound, nil, "not found")
		return
	}
	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/resources/"), "/exposure")
	snap, err := s.Store.LoadCurrent(r.Context())
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	out, err := analysis.New(graph.FromSnapshot(snap)).ExplainResourceExposure(id)
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	s.write(w, http.StatusOK, out)
}

func (s *Server) roleByName(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/roles/")
	snap, err := s.Store.LoadCurrent(r.Context())
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	for _, e := range snap.Entities {
		if e.Type == model.EntityRole && e.Name == name {
			s.write(w, http.StatusOK, e)
			return
		}
	}
	s.write(w, http.StatusNotFound, nil, "role not found")
}

func (s *Server) risks(w http.ResponseWriter, r *http.Request) {
	snap, err := s.Store.LoadCurrent(r.Context())
	if err != nil {
		s.write(w, http.StatusNotFound, nil, err.Error())
		return
	}
	s.write(w, http.StatusOK, analysis.New(graph.FromSnapshot(snap)).AnalyzeRisks())
}

type diffRequest struct {
	OldID string `json:"old_id"`
	NewID string `json:"new_id"`
}

func (s *Server) diff(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.write(w, http.StatusMethodNotAllowed, nil, "method not allowed")
		return
	}
	var req diffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.write(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	oldSnap, err := s.Store.Load(r.Context(), req.OldID)
	if err != nil {
		s.write(w, http.StatusNotFound, nil, "old snapshot: "+err.Error())
		return
	}
	newSnap, err := s.Store.Load(r.Context(), req.NewID)
	if err != nil {
		s.write(w, http.StatusNotFound, nil, "new snapshot: "+err.Error())
		return
	}
	s.write(w, http.StatusOK, diff.Diff(oldSnap, newSnap))
}
