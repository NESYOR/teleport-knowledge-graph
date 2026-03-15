package analysis

import (
	"fmt"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// ExplainResourceExposure identifies users and roles that can reach a resource.
func (s *Service) ExplainResourceExposure(resourceID string) (ExposureReport, error) {
	if s == nil || s.Graph == nil {
		return ExposureReport{}, fmt.Errorf("analysis service graph is nil")
	}
	if _, ok := s.Graph.Node(resourceID); !ok {
		return ExposureReport{}, fmt.Errorf("resource %q not found", resourceID)
	}

	nodes, _ := s.Graph.Export()
	users := make([]string, 0)
	roles := make([]string, 0)
	for _, n := range nodes {
		switch n.Type {
		case model.EntityUser:
			path := s.Graph.ShortestPath(n.ID, resourceID, map[model.RelationshipType]bool{
				model.RelUserHasRole:          true,
				model.RelUserCanAccessNode:    true,
				model.RelUserCanAccessDB:      true,
				model.RelUserCanAccessKube:    true,
				model.RelUserCanAccessDesktop: true,
				model.RelAccessPath:           true,
			})
			if len(path) > 0 {
				users = append(users, n.Name)
			}
		case model.EntityRole:
			path := s.Graph.ShortestPath(n.ID, resourceID, map[model.RelationshipType]bool{
				model.RelAccessPath: true,
			})
			if len(path) > 0 {
				roles = append(roles, n.Name)
			}
		}
	}
	return ExposureReport{ResourceID: resourceID, Users: sortStrings(users), Roles: sortStrings(roles)}, nil
}
