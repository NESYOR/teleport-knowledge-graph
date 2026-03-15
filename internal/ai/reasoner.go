package ai

import (
	"fmt"
	"strings"

	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// LocalRuleReasoner provides deterministic path-based explanations.
type LocalRuleReasoner struct {
	Graph *graph.Graph
}

// NewLocalRuleReasoner creates a deterministic rule reasoner.
func NewLocalRuleReasoner(g *graph.Graph) *LocalRuleReasoner {
	return &LocalRuleReasoner{Graph: g}
}

// WhyUserCanAccess explains whether and why a user can access a resource.
func (r *LocalRuleReasoner) WhyUserCanAccess(userID, resourceID string) (ReasoningResult, error) {
	if r == nil || r.Graph == nil {
		return ReasoningResult{}, fmt.Errorf("reasoner graph is nil")
	}
	path := r.Graph.ShortestPath(userID, resourceID, map[model.RelationshipType]bool{
		model.RelUserHasRole:          true,
		model.RelAccessPath:           true,
		model.RelUserCanAccessNode:    true,
		model.RelUserCanAccessDB:      true,
		model.RelUserCanAccessKube:    true,
		model.RelUserCanAccessDesktop: true,
	})
	q := fmt.Sprintf("Why can %s access %s?", userID, resourceID)
	if len(path) == 0 {
		return ReasoningResult{Question: q, Conclusion: "no_access_path", Path: nil, Evidence: []string{"no path found with allowed access edges"}}, nil
	}
	return ReasoningResult{Question: q, Conclusion: "access_path_found", Path: path, Evidence: []string{"shortest_path=" + strings.Join(path, "->")}}, nil
}
