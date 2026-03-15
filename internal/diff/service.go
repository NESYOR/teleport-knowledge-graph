package diff

import (
	"sort"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

// ChangeType describes snapshot-level change type.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// Change describes a single snapshot difference.
type Change struct {
	Type       ChangeType `json:"type"`
	EntityID   string     `json:"entity_id"`
	EntityType string     `json:"entity_type"`
	Severity   string     `json:"severity"`
	Reason     string     `json:"reason"`
}

// Result captures snapshot diff output.
type Result struct {
	Changes []Change `json:"changes"`
}

// Diff compares two snapshots and returns entity-level changes.
func Diff(oldSnap, newSnap model.Snapshot) Result {
	oldSet := make(map[string]model.Entity, len(oldSnap.Entities))
	newSet := make(map[string]model.Entity, len(newSnap.Entities))
	for _, e := range oldSnap.Entities {
		oldSet[e.ID] = e
	}
	for _, e := range newSnap.Entities {
		newSet[e.ID] = e
	}

	changes := make([]Change, 0)
	for id, n := range newSet {
		o, ok := oldSet[id]
		if !ok {
			changes = append(changes, Change{Type: ChangeAdded, EntityID: id, EntityType: string(n.Type), Severity: severity(ChangeAdded, n.Type), Reason: "entity added"})
			continue
		}
		if o.Name != n.Name || o.Type != n.Type || o.Namespace != n.Namespace || len(o.Labels) != len(n.Labels) {
			changes = append(changes, Change{Type: ChangeModified, EntityID: id, EntityType: string(n.Type), Severity: severity(ChangeModified, n.Type), Reason: "entity metadata changed"})
		}
	}
	for id, o := range oldSet {
		if _, ok := newSet[id]; !ok {
			changes = append(changes, Change{Type: ChangeRemoved, EntityID: id, EntityType: string(o.Type), Severity: severity(ChangeRemoved, o.Type), Reason: "entity removed"})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Severity == changes[j].Severity {
			if changes[i].Type == changes[j].Type {
				return changes[i].EntityID < changes[j].EntityID
			}
			return changes[i].Type < changes[j].Type
		}
		return severityRank(changes[i].Severity) > severityRank(changes[j].Severity)
	})
	return Result{Changes: changes}
}

func severity(kind ChangeType, entityType model.EntityType) string {
	if kind == ChangeRemoved {
		return "high"
	}
	switch entityType {
	case model.EntityRole, model.EntityUser:
		if kind == ChangeModified {
			return "high"
		}
		return "medium"
	case model.EntityNode, model.EntityDatabase, model.EntityKubernetesCluster, model.EntityApplication:
		return "medium"
	default:
		return "low"
	}
}

func severityRank(s string) int {
	switch s {
	case "high":
		return 3
	case "medium":
		return 2
	default:
		return 1
	}
}
