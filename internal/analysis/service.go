package analysis

import (
	"sort"

	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// Service provides analysis operations over a graph.
type Service struct {
	Graph *graph.Graph
}

// New creates an analysis service.
func New(g *graph.Graph) *Service {
	return &Service{Graph: g}
}

// UserAccess summarizes resources reachable by a user with evidence paths.
type UserAccess struct {
	UserID    string          `json:"user_id"`
	Resources []AccessBinding `json:"resources"`
}

// AccessBinding binds a target resource with path evidence.
type AccessBinding struct {
	ResourceID   string   `json:"resource_id"`
	ResourceType string   `json:"resource_type"`
	Path         []string `json:"path"`
}

// ExposureReport lists identities/roles that expose a resource.
type ExposureReport struct {
	ResourceID string   `json:"resource_id"`
	Users      []string `json:"users"`
	Roles      []string `json:"roles"`
}

// RiskFinding describes a risk signal with severity and evidence.
type RiskFinding struct {
	ID       string   `json:"id"`
	Severity string   `json:"severity"`
	Title    string   `json:"title"`
	Evidence []string `json:"evidence"`
}

// TopologySummary provides high-level inventory by type.
type TopologySummary struct {
	EntityCounts map[model.EntityType]int `json:"entity_counts"`
	TotalEdges   int                      `json:"total_edges"`
}

func resourceEntity(t model.EntityType) bool {
	switch t {
	case model.EntityNode, model.EntityDatabase, model.EntityKubernetesCluster, model.EntityWindowsDesktop, model.EntityApplication:
		return true
	default:
		return false
	}
}

func sortStrings(in []string) []string {
	out := append([]string(nil), in...)
	sort.Strings(out)
	return out
}
