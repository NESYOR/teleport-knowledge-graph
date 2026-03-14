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

// Load loads config from a path. For Phase 2 scaffolding, defaults are used
// when no file path is provided.
func Load(path string) (Config, error) {
	cfg := Config{
		Mode: "cli",
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
		Collection: CollectionConfig{
			Timeout: 2 * time.Minute,
		},
		Storage: StorageConfig{Path: "./snapshots"},
	}

	if path == "" {
		return cfg, nil
	}

	if _, err := os.Stat(path); err != nil {
		return Config{}, fmt.Errorf("config path %q not accessible: %w", path, err)
	}

	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
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
	return nil
}
