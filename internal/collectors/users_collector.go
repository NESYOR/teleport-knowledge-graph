package collectors

import (
	"context"

	"github.com/example/teleport-cluster-digital-twin/internal/model"
	"github.com/example/teleport-cluster-digital-twin/internal/teleport"
)

// UsersCollector collects users and USER_HAS_ROLE relationships.
type UsersCollector struct{}

// Name returns collector name.
func (UsersCollector) Name() string { return "users" }

// Collect gathers users and user-role edges.
func (UsersCollector) Collect(ctx context.Context, client teleport.ResourceClient) Result {
	users, err := client.ListUsers(ctx)
	if err != nil {
		return Result{Errors: []model.CollectorError{{Collector: "users", Kind: teleport.ErrorKind(err), Message: err.Error()}}}
	}

	entities := make([]model.Entity, 0, len(users))
	rels := make([]model.Relationship, 0)
	for _, u := range users {
		uid := model.DeterministicID("user", u.Name)
		entities = append(entities, model.Entity{ID: uid, Type: model.EntityUser, Name: u.Name, Cluster: ""})
		for _, role := range u.Roles {
			rid := model.DeterministicID("role", role)
			rels = append(rels, model.Relationship{
				ID:         model.DeterministicID("edge", string(model.RelUserHasRole), uid, rid),
				Type:       model.RelUserHasRole,
				FromID:     uid,
				ToID:       rid,
				Evidence:   "teleport.user.roles",
				Confidence: 1.0,
			})
		}
	}
	return Result{Entities: entities, Relationships: rels, Counts: map[string]int{"users": len(users)}}
}
