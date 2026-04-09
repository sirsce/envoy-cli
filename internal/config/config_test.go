package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envoy-cli/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if !cfg.Encrypt {
		t.Error("expected Encrypt to be true by default")
	}
	if len(cfg.EnvFiles) != 1 || cfg.EnvFiles[0] != ".env" {
		t.Errorf("unexpected default env files: %v", cfg.EnvFiles)
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envoy.json")

	cfg := &config.Config{
		Remote:   "s3://my-bucket/envs",
		Encrypt:  true,
		KeyFile:  ".envoy.key",
		EnvFiles: []string{".env", ".env.local"},
		Aliases:  map[string]string{"prod": "production"},
	}

	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Remote != cfg.Remote {
		t.Errorf("Remote mismatch: got %q want %q", loaded.Remote, cfg.Remote)
	}
	if loaded.KeyFile != cfg.KeyFile {
		t.Errorf("KeyFile mismatch: got %q want %q", loaded.KeyFile, cfg.KeyFile)
	}
	if len(loaded.EnvFiles) != 2 {
		t.Errorf("EnvFiles length mismatch: got %d", len(loaded.EnvFiles))
	}
	if loaded.Aliases["prod"] != "production" {
		t.Errorf("Alias not preserved")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	cfg, err := config.Load("/nonexistent/path/.envoy.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}
}

func TestValidate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cfg := &config.Config{
			Encrypt:  true,
			KeyFile:  ".envoy.key",
			EnvFiles: []string{".env"},
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing key file", func(t *testing.T) {
		cfg := &config.Config{
			Encrypt:  true,
			EnvFiles: []string{".env"},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing key_file")
		}
	})

	t.Run("no env files", func(t *testing.T) {
		cfg := &config.Config{
			Encrypt:  false,
			EnvFiles: []string{},
		}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for empty env files")
		}
	})
}

func TestSave_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envoy.json")

	cfg := config.DefaultConfig()
	cfg.KeyFile = ".envoy.key"
	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected file permission 0600, got %o", perm)
	}
}
