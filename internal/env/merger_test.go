package env

import (
	"testing"
)

func makeEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestMerger_KeepBase_NoConflict(t *testing.T) {
	m := NewMerger(MergeStrategyKeepBase)
	base := makeEntries("A", "1", "B", "2")
	override := makeEntries("C", "3")
	res, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Errorf("expected C to be added, got %v", res.Added)
	}
}

func TestMerger_KeepBase_Conflict(t *testing.T) {
	m := NewMerger(MergeStrategyKeepBase)
	base := makeEntries("A", "original")
	override := makeEntries("A", "new")
	res, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}
	if res.Entries[0].Value != "original" {
		t.Errorf("expected base value to be kept, got %q", res.Entries[0].Value)
	}
}

func TestMerger_Override_Conflict(t *testing.T) {
	m := NewMerger(MergeStrategyOverride)
	base := makeEntries("A", "original", "B", "2")
	override := makeEntries("A", "new")
	res, err := m.Merge(base, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overridden) != 1 || res.Overridden[0] != "A" {
		t.Errorf("expected A to be overridden, got %v", res.Overridden)
	}
	if res.Entries[0].Value != "new" {
		t.Errorf("expected overridden value, got %q", res.Entries[0].Value)
	}
}

func TestMerger_Error_Conflict(t *testing.T) {
	m := NewMerger(MergeStrategyError)
	base := makeEntries("A", "1")
	override := makeEntries("A", "2")
	_, err := m.Merge(base, override)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	conflict, ok := err.(*MergeConflictError)
	if !ok {
		t.Fatalf("expected MergeConflictError, got %T", err)
	}
	if conflict.Key != "A" {
		t.Errorf("expected conflict on key A, got %q", conflict.Key)
	}
}

func TestMerger_EmptyInputs(t *testing.T) {
	m := NewMerger(MergeStrategyOverride)
	res, err := m.Merge(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected empty result, got %d entries", len(res.Entries))
	}
}
