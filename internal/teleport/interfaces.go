package teleport

import "context"

// RoleRecord is a normalized role representation returned by Teleport client wrappers.
type RoleRecord struct {
	Name   string
	Logins []string
}

// UserRecord is a normalized user representation returned by Teleport client wrappers.
type UserRecord struct {
	Name  string
	Roles []string
}

// NodeRecord is a normalized node representation returned by Teleport client wrappers.
type NodeRecord struct {
	Name      string
	Namespace string
	Labels    map[string]string
}

// ResourceClient defines Teleport data operations required by collectors.
type ResourceClient interface {
	ClusterName(ctx context.Context) (string, error)
	CurrentUser(ctx context.Context) (string, error)
	ListRoles(ctx context.Context) ([]RoleRecord, error)
	ListUsers(ctx context.Context) ([]UserRecord, error)
	ListNodes(ctx context.Context) ([]NodeRecord, error)
}

// ClientFactory creates ResourceClient instances from runtime config.
type ClientFactory interface {
	New(ctx context.Context) (ResourceClient, AuthSource, error)
}
