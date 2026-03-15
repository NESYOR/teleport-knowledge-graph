package teleport

import (
	"context"
	"log/slog"

	"github.com/example/teleport-cluster-digital-twin/internal/config"
)

// Factory creates ResourceClient instances from runtime configuration.
type Factory struct {
	cfg    config.TeleportConfig
	logger *slog.Logger
}

// NewFactory creates a Teleport client factory.
func NewFactory(cfg config.TeleportConfig, logger *slog.Logger) *Factory {
	if logger == nil {
		logger = slog.Default()
	}
	return &Factory{cfg: cfg, logger: logger}
}

// New constructs a ResourceClient and identifies which auth source was used.
func (f *Factory) New(ctx context.Context) (ResourceClient, AuthSource, error) {
	_ = ctx
	var (
		auth AuthSource
		err  error
	)

	switch f.cfg.AuthMode {
	case "identity":
		auth, err = LoadIdentityAuth(f.cfg.IdentityFile)
	default:
		auth, err = LoadProfileAuth(f.cfg.ProfilePath)
	}
	if err != nil {
		return nil, AuthSource{}, err
	}

	// Phase 3 returns a deterministic in-process client implementation.
	// A real API-backed implementation can replace this type without changing collector contracts.
	client := &LocalStaticClient{cluster: f.cfg.Cluster}
	f.logger.Info("initialized teleport resource client", "auth_mode", auth.Mode, "auth_detail", auth.Detail)
	return client, auth, nil
}

// LocalStaticClient is a deterministic fallback client for development and tests.
type LocalStaticClient struct {
	cluster string
}

// ClusterName returns a configured or default cluster name.
func (c *LocalStaticClient) ClusterName(context.Context) (string, error) {
	if c.cluster == "" {
		return "example-cluster", nil
	}
	return c.cluster, nil
}

// CurrentUser returns a deterministic current user value.
func (c *LocalStaticClient) CurrentUser(context.Context) (string, error) {
	return "local-user", nil
}

// ListRoles returns sample visible roles.
func (c *LocalStaticClient) ListRoles(context.Context) ([]RoleRecord, error) {
	return []RoleRecord{
		{Name: "access", Logins: []string{"ubuntu"}},
		{Name: "prod-admin", Logins: []string{"root", "ubuntu"}},
	}, nil
}

// ListUsers returns sample visible users.
func (c *LocalStaticClient) ListUsers(context.Context) ([]UserRecord, error) {
	return []UserRecord{{Name: "alice", Roles: []string{"access"}}, {Name: "bob", Roles: []string{"prod-admin"}}}, nil
}

// ListNodes returns sample visible nodes.
func (c *LocalStaticClient) ListNodes(context.Context) ([]NodeRecord, error) {
	return []NodeRecord{
		{Name: "node-1", Namespace: "default", Labels: map[string]string{"env": "dev"}},
		{Name: "node-2", Namespace: "default", Labels: map[string]string{"env": "prod"}},
	}, nil
}
