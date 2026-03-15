package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func defaults() Config {
	return Config{
		Mode:       "cli",
		Log:        LogConfig{Level: "info", Format: "json"},
		Collection: CollectionConfig{Timeout: 2 * time.Minute, EnabledCollectors: []string{"cluster", "users", "roles", "nodes"}},
		Storage:    StorageConfig{Path: "./snapshots"},
		Teleport:   TeleportConfig{AuthMode: "profile"},
		API:        APIConfig{ListenAddr: ":8080"},
	}
}

// Load loads config from a simple YAML-like file path, overlaying defaults.
func Load(path string) (Config, error) {
	cfg := defaults()
	if path == "" {
		return cfg, Validate(cfg)
	}
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config path %q: %w", path, err)
	}
	defer f.Close()

	section := ""
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
			section = strings.TrimSuffix(line, ":")
			continue
		}
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		key := strings.TrimSpace(k)
		val := strings.Trim(strings.TrimSpace(v), "\"'")
		switch section + "." + key {
		case ".mode":
			cfg.Mode = val
		case "log.level":
			cfg.Log.Level = val
		case "log.format":
			cfg.Log.Format = val
		case "collection.timeout":
			if d, err := time.ParseDuration(val); err == nil {
				cfg.Collection.Timeout = d
			}
		case "collection.enabled_collectors":
			cfg.Collection.EnabledCollectors = parseList(val)
		case "storage.path":
			cfg.Storage.Path = val
		case "teleport.auth_mode":
			cfg.Teleport.AuthMode = val
		case "teleport.profile_path":
			cfg.Teleport.ProfilePath = val
		case "teleport.identity_file":
			cfg.Teleport.IdentityFile = val
		case "teleport.cluster":
			cfg.Teleport.Cluster = val
		case "teleport.proxy_addr":
			cfg.Teleport.ProxyAddr = val
		case "teleport.insecure":
			b, _ := strconv.ParseBool(val)
			cfg.Teleport.Insecure = b
		case "api.listen_addr":
			cfg.API.ListenAddr = val
		}
	}
	if err := s.Err(); err != nil {
		return Config{}, fmt.Errorf("scan config file: %w", err)
	}
	return cfg, Validate(cfg)
}

func parseList(v string) []string {
	trim := strings.TrimSpace(v)
	trim = strings.TrimPrefix(trim, "[")
	trim = strings.TrimSuffix(trim, "]")
	if trim == "" {
		return nil
	}
	parts := strings.Split(trim, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.Trim(strings.TrimSpace(p), "\"'")
		if t != "" {
			out = append(out, t)
		}
	}
	return out
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
