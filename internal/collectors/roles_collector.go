package collectors

import (
	"context"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// RolesCollector collects roles and role-login relationships.
type RolesCollector struct{}

// Name returns collector name.
func (RolesCollector) Name() string { return "roles" }

// Collect gathers role entities and login entities.
func (RolesCollector) Collect(ctx context.Context, client teleport.ResourceClient) Result {
	roles, err := client.ListRoles(ctx)
	if err != nil {
		return Result{Errors: []model.CollectorError{{Collector: "roles", Kind: teleport.ErrorKind(err), Message: err.Error()}}}
	}

	entities := make([]model.Entity, 0, len(roles))
	rels := make([]model.Relationship, 0)
	loginCount := 0
	for _, r := range roles {
		rid := model.DeterministicID("role", r.Name)
		entities = append(entities, model.Entity{ID: rid, Type: model.EntityRole, Name: r.Name})
		for _, login := range r.Logins {
			lid := model.DeterministicID("login", login)
			entities = append(entities, model.Entity{ID: lid, Type: model.EntityLogin, Name: login})
			rels = append(rels, model.Relationship{
				ID:         model.DeterministicID("edge", string(model.RelRoleAllowsLogin), rid, lid),
				Type:       model.RelRoleAllowsLogin,
				FromID:     rid,
				ToID:       lid,
				Evidence:   "teleport.role.logins",
				Confidence: 1.0,
			})
			loginCount++
		}
	}

	return Result{Entities: entities, Relationships: rels, Counts: map[string]int{"roles": len(roles), "logins": loginCount}}
}
