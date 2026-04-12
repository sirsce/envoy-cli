package env

import (
	"testing"
	"time"
)

var fixedSnapTime = time.Date(2024, 1, 15, 12, 0, 0, 42, time.UTC)

func newTestSnapshotManager() *SnapshotManager {
	sm := NewSnapshotManager()
	sm.clock = func() time.Time { return fixedSnapTime }
	return sm
}

func TestSnapshotManager_TakeAndGet(t *testing.T) {
	sm := newTestSnapshotManager()
	entries := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}

	snap := sm.Take("test", entries)
	if snap.Label != "test" {
		t.Errorf("expected label 'test', got %q", snap.Label)
	}
	if len(snap.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(snap.Entries))
	}

	got, ok := sm.Get(snap.ID)
	if !ok {
		t.Fatal("expected snapshot to be found")
	}
	if got.ID != snap.ID {
		t.Errorf("ID mismatch: %q vs %q", got.ID, snap.ID)
	}
}

func TestSnapshotManager_DoesNotMutateOriginal(t *testing.T) {
	sm := newTestSnapshotManager()
	entries := []Entry{{Key: "KEY", Value: "original"}}
	snap := sm.Take("immutable", entries)

	entries[0].Value = "mutated"
	if snap.Entries[0].Value != "original" {
		t.Error("snapshot was mutated by original slice change")
	}
}

func TestSnapshotManager_List(t *testing.T) {
	sm := newTestSnapshotManager()
	sm.Take("alpha", []Entry{{Key: "A", Value: "1"}})
	sm.Take("beta", []Entry{{Key: "B", Value: "2"}})

	list := sm.List()
	if len(list) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(list))
	}
}

func TestSnapshotManager_Delete(t *testing.T) {
	sm := newTestSnapshotManager()
	snap := sm.Take("to-delete", []Entry{{Key: "X", Value: "y"}})

	if !sm.Delete(snap.ID) {
		t.Error("expected Delete to return true")
	}
	if sm.Delete(snap.ID) {
		t.Error("expected second Delete to return false")
	}
	if _, ok := sm.Get(snap.ID); ok {
		t.Error("expected snapshot to be gone after delete")
	}
}

func TestSnapshotManager_Restore(t *testing.T) {
	sm := newTestSnapshotManager()
	entries := []Entry{{Key: "ENV", Value: "prod"}}
	snap := sm.Take("restore-test", entries)

	restored, err := sm.Restore(snap.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(restored) != 1 || restored[0].Key != "ENV" {
		t.Errorf("unexpected restored entries: %+v", restored)
	}
}

func TestSnapshotManager_Restore_NotFound(t *testing.T) {
	sm := newTestSnapshotManager()
	_, err := sm.Restore("nonexistent-id")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}
