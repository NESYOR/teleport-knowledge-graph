package teleport

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/teleport-cluster-digital-twin/internal/config"
)

func TestFactoryNewProfileAuth(t *testing.T) {
	dir := t.TempDir()
	cfg := config.TeleportConfig{AuthMode: "profile", ProfilePath: dir, Cluster: "c1"}
	f := NewFactory(cfg, slog.Default())
	client, auth, err := f.New(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.Mode != "profile" {
		t.Fatalf("expected profile mode, got %s", auth.Mode)
	}
	name, err := client.ClusterName(context.Background())
	if err != nil || name != "c1" {
		t.Fatalf("unexpected cluster name: %s err=%v", name, err)
	}
}

func TestFactoryNewIdentityAuthMissing(t *testing.T) {
	cfg := config.TeleportConfig{AuthMode: "identity", IdentityFile: ""}
	f := NewFactory(cfg, slog.Default())
	_, _, err := f.New(context.Background())
	if err == nil {
		t.Fatal("expected error for missing identity file")
	}
}

func TestLoadIdentityAuth(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "identity")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	auth, err := LoadIdentityAuth(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.Mode != "identity" {
		t.Fatalf("expected identity mode, got %s", auth.Mode)
	}
}
