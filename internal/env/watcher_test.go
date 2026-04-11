package env

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempEnv(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return path
}

func TestWatcher_NoChangeAfterSeed(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, "KEY=value\n")

	w := NewWatcher(path, 20*time.Millisecond)
	if err := w.Seed(); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	ch := w.Start()
	defer w.Stop()

	select {
	case ev := <-ch:
		if ev.Changed {
			t.Errorf("unexpected change event after seed")
		}
	case <-time.After(80 * time.Millisecond):
		// expected: no change
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := writeTempEnv(t, dir, "KEY=original\n")

	w := NewWatcher(path, 20*time.Millisecond)
	if err := w.Seed(); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	ch := w.Start()
	defer w.Stop()

	// Modify the file after a short delay.
	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=modified\n"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Err != nil {
			t.Fatalf("unexpected error: %v", ev.Err)
		}
		if !ev.Changed {
			t.Errorf("expected Changed=true")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("timed out waiting for change event")
	}
}

func TestWatcher_MissingFile(t *testing.T) {
	w := NewWatcher("/nonexistent/.env", 20*time.Millisecond)
	ch := w.Start()
	defer w.Stop()

	select {
	case ev := <-ch:
		if ev.Err == nil {
			t.Error("expected error for missing file")
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("timed out waiting for error event")
	}
}

func TestWatcher_SeedMissingFile(t *testing.T) {
	w := NewWatcher("/nonexistent/.env", 20*time.Millisecond)
	if err := w.Seed(); err == nil {
		t.Error("expected Seed to return error for missing file")
	}
}
