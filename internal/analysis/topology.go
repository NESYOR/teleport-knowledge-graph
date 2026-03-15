package analysis

import "github.com/example/teleport-cluster-digital-twin/internal/model"

// Topology returns counts by entity type and total relationship count.
func (s *Service) Topology() TopologySummary {
	if s == nil || s.Graph == nil {
		return TopologySummary{EntityCounts: map[model.EntityType]int{}}
	}
	nodes, edges := s.Graph.Export()
	counts := map[model.EntityType]int{}
	for _, n := range nodes {
		counts[n.Type]++
	}
	return TopologySummary{EntityCounts: counts, TotalEdges: len(edges)}
}
