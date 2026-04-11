package env_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envoy-cli/internal/env"
)

// TestWatcher_Integration_MultipleChanges verifies that successive writes each
// produce a distinct change event.
func TestWatcher_Integration_MultipleChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")

	writeFile := func(content string) {
		t.Helper()
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	writeFile("A=1\n")

	w := env.NewWatcher(path, 15*time.Millisecond)
	if err := w.Seed(); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	ch := w.Start()
	defer w.Stop()

	changes := 0
	done := make(chan struct{})

	go func() {
		defer close(done)
		for ev := range ch {
			if ev.Err != nil {
				t.Errorf("unexpected error: %v", ev.Err)
				return
			}
			if ev.Changed {
				changes++
				if changes == 2 {
					return
				}
			}
		}
	}()

	time.Sleep(30 * time.Millisecond)
	writeFile("A=2\n")
	time.Sleep(60 * time.Millisecond)
	writeFile("A=3\n")

	select {
	case <-done:
		if changes < 2 {
			t.Errorf("expected at least 2 change events, got %d", changes)
		}
	case <-time.After(500 * time.Millisecond):
		w.Stop()
		t.Errorf("timed out: only %d change(s) detected", changes)
	}
}
