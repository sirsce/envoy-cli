package storage_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/storage"
)

func TestLocalBackend_WriteRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", ".env")

	b := storage.NewLocalBackend(path)

	exists, err := b.Exists()
	if err != nil {
		t.Fatalf("Exists error: %v", err)
	}
	if exists {
		t.Fatal("expected file to not exist yet")
	}

	const content = "KEY=value\nFOO=bar\n"
	if err := b.Write(strings.NewReader(content)); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	exists, err = b.Exists()
	if err != nil {
		t.Fatalf("Exists after write error: %v", err)
	}
	if !exists {
		t.Fatal("expected file to exist after write")
	}

	rc, err := b.Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	defer rc.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if string(data) != content {
		t.Errorf("content mismatch: got %q want %q", string(data), content)
	}
}

func TestLocalBackend_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	b := storage.NewLocalBackend(path)
	if err := b.Write(strings.NewReader("SECRET=abc")); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("expected permissions 0600, got %o", perm)
	}
}
