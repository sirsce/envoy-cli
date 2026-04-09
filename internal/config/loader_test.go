package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/config"
)

func TestFindConfigFile_Found(t *testing.T) {
	root := t.TempDir()
	cfgPath := filepath.Join(root, ".envoy.json")

	cfg := config.DefaultConfig()
	cfg.KeyFile = ".envoy.key"
	if err := config.Save(cfg, cfgPath); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Search from a subdirectory.
	sub := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	found := config.FindConfigFile(sub)
	if found != cfgPath {
		t.Errorf("FindConfigFile: got %q, want %q", found, cfgPath)
	}
}

func TestFindConfigFile_NotFound(t *testing.T) {
	dir := t.TempDir()
	found := config.FindConfigFile(dir)
	if found != "" {
		t.Errorf("expected empty string, got %q", found)
	}
}

func TestMustLoad_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected MustLoad to panic on invalid JSON")
		}
	}()

	f, err := os.CreateTemp("", "bad-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("{invalid json")
	f.Close()

	config.MustLoad(f.Name())
}

func TestMustLoad_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envoy.json")

	cfg := &config.Config{
		Remote:   "local:///tmp/envs",
		Encrypt:  false,
		EnvFiles: []string{".env"},
		Aliases:  map[string]string{},
	}
	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded := config.MustLoad(path)
	if loaded.Remote != cfg.Remote {
		t.Errorf("Remote mismatch: got %q want %q", loaded.Remote, cfg.Remote)
	}
}
