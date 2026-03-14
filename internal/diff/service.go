package diff

import "github.com/example/teleport-cluster-digital-twin/internal/model"

// ChangeType describes snapshot-level change type.
type ChangeType string

const (
	ChangeAdded    ChangeType = "added"
	ChangeRemoved  ChangeType = "removed"
	ChangeModified ChangeType = "modified"
)

// Change describes a single snapshot difference.
type Change struct {
	Type     ChangeType `json:"type"`
	EntityID string     `json:"entity_id"`
	Severity string     `json:"severity"`
}

// Diff compares two snapshots and returns coarse entity changes.
func Diff(oldSnap, newSnap model.Snapshot) []Change {
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
			changes = append(changes, Change{Type: ChangeAdded, EntityID: id, Severity: "low"})
			continue
		}
		if o.Name != n.Name || o.Type != n.Type {
			changes = append(changes, Change{Type: ChangeModified, EntityID: id, Severity: "medium"})
		}
	}
	for id := range oldSet {
		if _, ok := newSet[id]; !ok {
			changes = append(changes, Change{Type: ChangeRemoved, EntityID: id, Severity: "high"})
		}
	}
	return changes
}
