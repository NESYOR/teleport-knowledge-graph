package ai

import (
	"fmt"
	"sort"

	"github.com/example/teleport-cluster-digital-twin/internal/graph"
	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// GraphSummarizer emits deterministic summaries from graph data.
type GraphSummarizer struct {
	Graph *graph.Graph
}

// NewGraphSummarizer creates a deterministic graph summarizer.
func NewGraphSummarizer(g *graph.Graph) *GraphSummarizer {
	return &GraphSummarizer{Graph: g}
}

// SummarizeUser summarizes user role memberships and immediate reachable targets.
func (s *GraphSummarizer) SummarizeUser(userID string, maxResources int) (UserSummary, error) {
	if s == nil || s.Graph == nil {
		return UserSummary{}, fmt.Errorf("summarizer graph is nil")
	}
	user, ok := s.Graph.Node(userID)
	if !ok || user.Type != model.EntityUser {
		return UserSummary{}, fmt.Errorf("user %q not found", userID)
	}
	roles := make([]string, 0)
	targets := make([]string, 0)
	for _, e := range s.Graph.OutgoingEdges(userID) {
		if e.Type == model.RelUserHasRole {
			if role, ok := s.Graph.Node(e.ToID); ok {
				roles = append(roles, role.Name)
			}
		}
	}
	if maxResources <= 0 {
		maxResources = 10
	}
	for _, id := range s.Graph.Traverse(userID, 3, map[model.RelationshipType]bool{model.RelAccessPath: true, model.RelUserHasRole: true}) {
		if len(targets) >= maxResources {
			break
		}
		n, ok := s.Graph.Node(id)
		if !ok || !isResourceType(n.Type) {
			continue
		}
		targets = append(targets, n.ID)
	}
	sort.Strings(roles)
	sort.Strings(targets)
	return UserSummary{UserID: userID, Roles: roles, ReachableTargets: targets}, nil
}

// SummarizeRole summarizes role grants and login permissions.
func (s *GraphSummarizer) SummarizeRole(roleID string) (RoleSummary, error) {
	if s == nil || s.Graph == nil {
		return RoleSummary{}, fmt.Errorf("summarizer graph is nil")
	}
	role, ok := s.Graph.Node(roleID)
	if !ok || role.Type != model.EntityRole {
		return RoleSummary{}, fmt.Errorf("role %q not found", roleID)
	}
	logins := make([]string, 0)
	caps := make([]string, 0)
	for _, e := range s.Graph.OutgoingEdges(roleID) {
		switch e.Type {
		case model.RelRoleAllowsLogin:
			if login, ok := s.Graph.Node(e.ToID); ok {
				logins = append(logins, login.Name)
			}
		case model.RelRoleGrantsVerb, model.RelAccessPath:
			caps = append(caps, string(e.Type)+":"+e.ToID)
		}
	}
	sort.Strings(logins)
	sort.Strings(caps)
	return RoleSummary{RoleID: roleID, AllowedLogins: logins, Capabilities: caps}, nil
}

// SummarizeResource summarizes who can reach a resource and key labels.
func (s *GraphSummarizer) SummarizeResource(resourceID string) (ResourceSummary, error) {
	if s == nil || s.Graph == nil {
		return ResourceSummary{}, fmt.Errorf("summarizer graph is nil")
	}
	res, ok := s.Graph.Node(resourceID)
	if !ok {
		return ResourceSummary{}, fmt.Errorf("resource %q not found", resourceID)
	}
	exposed := make([]string, 0)
	for _, n := range s.Graph.NodesByType(model.EntityUser) {
		path := s.Graph.ShortestPath(n.ID, resourceID, map[model.RelationshipType]bool{model.RelUserHasRole: true, model.RelAccessPath: true, model.RelUserCanAccessNode: true, model.RelUserCanAccessDB: true, model.RelUserCanAccessKube: true, model.RelUserCanAccessDesktop: true})
		if len(path) > 0 {
			exposed = append(exposed, n.Name)
		}
	}
	labels := make([]string, 0, len(res.Labels))
	for k, v := range res.Labels {
		labels = append(labels, k+"="+v)
	}
	sort.Strings(exposed)
	sort.Strings(labels)
	return ResourceSummary{ResourceID: resourceID, ExposedTo: exposed, Labels: labels}, nil
}

func isResourceType(t model.EntityType) bool {
	switch t {
	case model.EntityNode, model.EntityDatabase, model.EntityKubernetesCluster, model.EntityWindowsDesktop, model.EntityApplication:
		return true
	default:
		return false
	}
}
