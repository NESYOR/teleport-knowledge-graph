package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Config contains top-level runtime configuration.
type Config struct {
	Mode       string
	Log        LogConfig
	Collection CollectionConfig
	Storage    StorageConfig
	Teleport   TeleportConfig
	API        APIConfig
}

// LogConfig controls structured logging behavior.
type LogConfig struct {
	Level  string
	Format string
}

// CollectionConfig controls collection timing and toggles.
type CollectionConfig struct {
	Timeout           time.Duration
	EnabledCollectors []string
}

// StorageConfig controls snapshot persistence location.
type StorageConfig struct {
	Path string
}

// TeleportConfig controls auth and target cluster information.
type TeleportConfig struct {
	ProxyAddr    string
	Cluster      string
	AuthMode     string
	ProfilePath  string
	IdentityFile string
	Insecure     bool
}

// APIConfig controls REST API server behavior.
type APIConfig struct {
	ListenAddr string
}

// Load loads config from a path. For Phase scaffolding, defaults are used
// when no file path is provided.
func Load(path string) (Config, error) {
	cfg := Config{
		Mode:       "cli",
		Log:        LogConfig{Level: "info", Format: "json"},
		Collection: CollectionConfig{Timeout: 2 * time.Minute, EnabledCollectors: []string{"cluster", "users", "roles", "nodes"}},
		Storage:    StorageConfig{Path: "./snapshots"},
		Teleport:   TeleportConfig{AuthMode: "profile"},
		API:        APIConfig{ListenAddr: ":8080"},
	}
	if path == "" {
		return cfg, Validate(cfg)
	}
	if _, err := os.Stat(path); err != nil {
		return Config{}, fmt.Errorf("config path %q not accessible: %w", path, err)
	}
	return cfg, Validate(cfg)
}

// Validate validates runtime configuration.
func Validate(cfg Config) error {
	if cfg.Mode == "" {
		return errors.New("mode is required")
	}
	if cfg.Collection.Timeout <= 0 {
		return errors.New("collection timeout must be > 0")
	}
	if cfg.Storage.Path == "" {
		return errors.New("storage path is required")
	}
	if cfg.Teleport.AuthMode == "" {
		return errors.New("teleport auth mode is required")
	}
	if cfg.API.ListenAddr == "" {
		return errors.New("api listen addr is required")
	}
	return nil
}
