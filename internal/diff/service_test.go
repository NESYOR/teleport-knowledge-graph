package diff

import (
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
)

func TestDiffClassifiesChanges(t *testing.T) {
	oldSnap := model.Snapshot{Entities: []model.Entity{
		{ID: "user:1", Type: model.EntityUser, Name: "alice"},
		{ID: "node:1", Type: model.EntityNode, Name: "n1"},
	}}
	newSnap := model.Snapshot{Entities: []model.Entity{
		{ID: "user:1", Type: model.EntityUser, Name: "alice2"}, // modified
		{ID: "role:1", Type: model.EntityRole, Name: "admin"},  // added
	}}

	res := Diff(oldSnap, newSnap)
	if len(res.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(res.Changes))
	}

	var sawHigh bool
	for _, c := range res.Changes {
		if c.Severity == "high" {
			sawHigh = true
		}
	}
	if !sawHigh {
		t.Fatal("expected at least one high severity change")
	}
}
