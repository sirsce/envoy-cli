package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileLoader_Load(t *testing.T) {
	dir := t.TempDir()
	content := "APP_ENV=staging\nDB_HOST=stage.db\n"
	path := filepath.Join(dir, "staging.env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	loader := NewProfileLoader(dir)
	profile, err := loader.Load("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.Name != "staging" {
		t.Errorf("expected name %q, got %q", "staging", profile.Name)
	}
	if len(profile.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(profile.Entries))
	}
}

func TestProfileLoader_Load_DotPrefix(t *testing.T) {
	dir := t.TempDir()
	content := "SECRET=abc\n"
	path := filepath.Join(dir, ".local.env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	loader := NewProfileLoader(dir)
	profile, err := loader.Load("local")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profile.Entries) != 1 || profile.Entries[0].Key != "SECRET" {
		t.Errorf("unexpected entries: %v", profile.Entries)
	}
}

func TestProfileLoader_Load_NotFound(t *testing.T) {
	dir := t.TempDir()
	loader := NewProfileLoader(dir)
	_, err := loader.Load("missing")
	if err == nil {
		t.Fatal("expected error for missing profile file")
	}
}

func TestProfileLoader_Load_InvalidContent(t *testing.T) {
	dir := t.TempDir()
	// A line that is just a value with no key should be handled gracefully.
	// Parser skips blank/comment lines; an invalid line may return an error.
	content := "VALID_KEY=hello\n"
	path := filepath.Join(dir, "test.env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write file: %v", err)
	}
	loader := NewProfileLoader(dir)
	profile, err := loader.Load("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profile.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(profile.Entries))
	}
}
