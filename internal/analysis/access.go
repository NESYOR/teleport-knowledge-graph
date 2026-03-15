package analysis

import (
	"fmt"
	"sort"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// ExplainUserAccess computes user-to-resource access paths through role membership.
func (s *Service) ExplainUserAccess(userName string) (UserAccess, error) {
	if s == nil || s.Graph == nil {
		return UserAccess{}, fmt.Errorf("analysis service graph is nil")
	}
	uid := model.DeterministicID("user", userName)
	if _, ok := s.Graph.Node(uid); !ok {
		return UserAccess{}, fmt.Errorf("user %q not found", userName)
	}

	nodes, _ := s.Graph.Export()
	bindings := make([]AccessBinding, 0)
	for _, n := range nodes {
		if !resourceEntity(n.Type) {
			continue
		}
		path := s.Graph.ShortestPath(uid, n.ID, map[model.RelationshipType]bool{
			model.RelUserHasRole:          true,
			model.RelUserCanAccessNode:    true,
			model.RelUserCanAccessDB:      true,
			model.RelUserCanAccessKube:    true,
			model.RelUserCanAccessDesktop: true,
			model.RelAccessPath:           true,
		})
		if len(path) == 0 {
			continue
		}
		bindings = append(bindings, AccessBinding{ResourceID: n.ID, ResourceType: string(n.Type), Path: path})
	}
	sort.Slice(bindings, func(i, j int) bool { return bindings[i].ResourceID < bindings[j].ResourceID })
	return UserAccess{UserID: uid, Resources: bindings}, nil
}
