package env

import (
	"testing"
)

func TestDiff_NoDifferences(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	diffs := Diff(base, target)
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestDiff_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	diffs := Diff(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Type != DiffAdded || diffs[0].Key != "NEW_KEY" || diffs[0].NewValue != "value" {
		t.Errorf("unexpected diff entry: %+v", diffs[0])
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "old"}
	target := map[string]string{"FOO": "bar"}

	diffs := Diff(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Type != DiffRemoved || diffs[0].Key != "OLD_KEY" || diffs[0].OldValue != "old" {
		t.Errorf("unexpected diff entry: %+v", diffs[0])
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	base := map[string]string{"FOO": "old_value"}
	target := map[string]string{"FOO": "new_value"}

	diffs := Diff(base, target)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if d.Type != DiffChanged || d.Key != "FOO" || d.OldValue != "old_value" || d.NewValue != "new_value" {
		t.Errorf("unexpected diff entry: %+v", d)
	}
}

func TestDiff_Mixed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "1", "B": "changed", "D": "4"}

	diffs := Diff(base, target)
	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(diffs))
	}
}

func TestHasDiff_True(t *testing.T) {
	base := map[string]string{"X": "1"}
	target := map[string]string{"X": "2"}
	if !HasDiff(base, target) {
		t.Error("expected HasDiff to return true")
	}
}

func TestHasDiff_False(t *testing.T) {
	base := map[string]string{"X": "1"}
	target := map[string]string{"X": "1"}
	if HasDiff(base, target) {
		t.Error("expected HasDiff to return false")
	}
}
