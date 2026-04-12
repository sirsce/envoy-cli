package env

import (
	"strings"
	"testing"
)

func TestCompareSnapshots_Identical(t *testing.T) {
	sm := newTestSnapshotManager()
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	snap1 := sm.Take("v1", entries)
	snap2 := sm.Take("v2", entries)

	result, err := CompareSnapshots(sm, snap1.ID, snap2.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if HasDiff(result.Diff) {
		t.Error("expected no diff between identical snapshots")
	}
	if !strings.Contains(result.String(), "identical") {
		t.Errorf("expected 'identical' in string: %s", result.String())
	}
}

func TestCompareSnapshots_WithChanges(t *testing.T) {
	sm := newTestSnapshotManager()
	v1 := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	v2 := []Entry{{Key: "A", Value: "updated"}, {Key: "C", Value: "3"}}

	snap1 := sm.Take("v1", v1)
	snap2 := sm.Take("v2", v2)

	result, err := CompareSnapshots(sm, snap1.ID, snap2.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !HasDiff(result.Diff) {
		t.Error("expected diff between different snapshots")
	}
	if len(result.Diff.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(result.Diff.Added))
	}
	if len(result.Diff.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Diff.Removed))
	}
	if len(result.Diff.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(result.Diff.Changed))
	}
	str := result.String()
	if !strings.Contains(str, "+1") || !strings.Contains(str, "-1") {
		t.Errorf("unexpected summary string: %s", str)
	}
}

func TestCompareSnapshots_FromNotFound(t *testing.T) {
	sm := newTestSnapshotManager()
	snap := sm.Take("v1", []Entry{{Key: "A", Value: "1"}})
	_, err := CompareSnapshots(sm, "missing-id", snap.ID)
	if err == nil {
		t.Error("expected error for missing from snapshot")
	}
}

func TestCompareSnapshots_ToNotFound(t *testing.T) {
	sm := newTestSnapshotManager()
	snap := sm.Take("v1", []Entry{{Key: "A", Value: "1"}})
	_, err := CompareSnapshots(sm, snap.ID, "missing-id")
	if err == nil {
		t.Error("expected error for missing to snapshot")
	}
}
