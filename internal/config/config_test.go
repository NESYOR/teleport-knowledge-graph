package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("load defaults: %v", err)
	}
	if cfg.API.ListenAddr == "" || cfg.Teleport.AuthMode == "" {
		t.Fatalf("expected defaults populated: %#v", cfg)
	}
}

func TestLoadYAML(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "config.yaml")
	content := []byte(`mode: serve
collection:
  timeout: 30s
storage:
  path: ./out
teleport:
  auth_mode: identity
  identity_file: /tmp/ident
api:
  listen_addr: :9090
`)
	if err := os.WriteFile(p, content, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("load yaml: %v", err)
	}
	if cfg.Mode != "serve" || cfg.API.ListenAddr != ":9090" || cfg.Teleport.AuthMode != "identity" {
		t.Fatalf("unexpected cfg: %#v", cfg)
	}
}
