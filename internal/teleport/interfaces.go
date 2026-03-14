package teleport

import "context"

// ResourceClient defines Teleport data operations required by collectors.
type ResourceClient interface {
	ClusterName(ctx context.Context) (string, error)
	CurrentUser(ctx context.Context) (string, error)
	ListRoles(ctx context.Context) ([]string, error)
	ListUsers(ctx context.Context) ([]string, error)
	ListNodes(ctx context.Context) ([]string, error)
}

// ClientFactory creates ResourceClient instances from runtime config.
type ClientFactory interface {
	New(ctx context.Context) (ResourceClient, error)
}
