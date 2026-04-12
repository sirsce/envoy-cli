package env

import (
	"strings"
	"testing"
)

// TestSnapshot_Integration_TakeRestoreAndDiff exercises the full snapshot
// lifecycle: parse → take snapshot → mutate → take snapshot → diff → restore.
func TestSnapshot_Integration_TakeRestoreAndDiff(t *testing.T) {
	raw := "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=development\n"
	p := NewParser()
	initial, err := p.ParseBytes([]byte(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	sm := NewSnapshotManager()
	snap1 := sm.Take("initial", initial)

	// Simulate a change: update DB_PORT, add SECRET, remove APP_ENV
	updated := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5433"},
		{Key: "SECRET", Value: "abc123"},
	}
	snap2 := sm.Take("updated", updated)

	result, err := CompareSnapshots(sm, snap1.ID, snap2.ID)
	if err != nil {
		t.Fatalf("compare error: %v", err)
	}
	if !HasDiff(result.Diff) {
		t.Fatal("expected diff between snapshots")
	}
	if len(result.Diff.Added) != 1 || result.Diff.Added[0].Key != "SECRET" {
		t.Errorf("expected SECRET added, got: %+v", result.Diff.Added)
	}
	if len(result.Diff.Removed) != 1 || result.Diff.Removed[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV removed, got: %+v", result.Diff.Removed)
	}
	if len(result.Diff.Changed) != 1 || result.Diff.Changed[0].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT changed, got: %+v", result.Diff.Changed)
	}

	// Restore original and verify
	restored, err := sm.Restore(snap1.ID)
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if len(restored) != 3 {
		t.Errorf("expected 3 restored entries, got %d", len(restored))
	}

	// Export restored to verify content
	ex := NewExporter(restored)
	out, err := ex.Format(FormatDotenv)
	if err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST in restored export, got: %s", out)
	}
	if !strings.Contains(out, "APP_ENV=development") {
		t.Errorf("expected APP_ENV in restored export, got: %s", out)
	}
}

func TestSnapshot_Integration_ListAndDelete(t *testing.T) {
	sm := NewSnapshotManager()
	entries := []Entry{{Key: "X", Value: "1"}}

	s1 := sm.Take("snap-a", entries)
	s2 := sm.Take("snap-b", entries)

	if len(sm.List()) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(sm.List()))
	}

	sm.Delete(s1.ID)
	if len(sm.List()) != 1 {
		t.Errorf("expected 1 snapshot after delete, got %d", len(sm.List()))
	}
	if _, ok := sm.Get(s2.ID); !ok {
		t.Error("expected snap-b to still exist")
	}
}
