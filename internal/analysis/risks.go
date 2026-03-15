package analysis

import (
	"sort"
	"strings"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// AnalyzeRisks runs deterministic rule-based risk checks over graph entities.
func (s *Service) AnalyzeRisks() []RiskFinding {
	if s == nil || s.Graph == nil {
		return nil
	}
	nodes, _ := s.Graph.Export()
	findings := make([]RiskFinding, 0)
	for _, n := range nodes {
		switch n.Type {
		case model.EntityRole:
			if strings.Contains(strings.ToLower(n.Name), "admin") {
				findings = append(findings, RiskFinding{
					ID:       model.DeterministicID("risk", "admin-role", n.ID),
					Severity: "high",
					Title:    "Privileged role detected",
					Evidence: []string{n.ID, n.Name},
				})
			}
		case model.EntityLogin:
			if n.Name == "root" {
				findings = append(findings, RiskFinding{
					ID:       model.DeterministicID("risk", "root-login", n.ID),
					Severity: "high",
					Title:    "Root login allowed",
					Evidence: []string{n.ID, n.Name},
				})
			}
		case model.EntityNode:
			for k, v := range n.Labels {
				if v == "*" {
					findings = append(findings, RiskFinding{
						ID:       model.DeterministicID("risk", "wildcard-label", n.ID, k),
						Severity: "medium",
						Title:    "Wildcard label detected",
						Evidence: []string{n.ID, k + "=*"},
					})
				}
			}
		}
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return findings
}
