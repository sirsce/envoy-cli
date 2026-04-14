package env

import (
	"testing"
)

func makeDedupEntries(pairs ...string) []Entry {
	entries := make([]Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDeduplicator_NoDuplicates(t *testing.T) {
	d := NewDeduplicator(DeduplicateKeepFirst)
	entries := makeDedupEntries("A", "1", "B", "2", "C", "3")
	res, err := d.Deduplicate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %d", len(res.Duplicates))
	}
}

func TestDeduplicator_KeepFirst(t *testing.T) {
	d := NewDeduplicator(DeduplicateKeepFirst)
	entries := makeDedupEntries("A", "first", "B", "only", "A", "second")
	res, err := d.Deduplicate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	if res.Entries[0].Value != "first" {
		t.Errorf("expected first value to be kept, got %q", res.Entries[0].Value)
	}
	if len(res.Duplicates) != 1 || res.Duplicates[0].Key != "A" {
		t.Errorf("expected duplicate record for A")
	}
}

func TestDeduplicator_KeepLast(t *testing.T) {
	d := NewDeduplicator(DeduplicateKeepLast)
	entries := makeDedupEntries("A", "first", "B", "only", "A", "second")
	res, err := d.Deduplicate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "second" {
		t.Errorf("expected last value to be kept, got %q", res.Entries[0].Value)
	}
}

func TestDeduplicator_ErrorStrategy(t *testing.T) {
	d := NewDeduplicator(DeduplicateError)
	entries := makeDedupEntries("X", "1", "X", "2")
	_, err := d.Deduplicate(entries)
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestDeduplicator_ErrorStrategy_NoDuplicates(t *testing.T) {
	d := NewDeduplicator(DeduplicateError)
	entries := makeDedupEntries("X", "1", "Y", "2")
	res, err := d.Deduplicate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
}

func TestHasDuplicates(t *testing.T) {
	if HasDuplicates(makeDedupEntries("A", "1", "B", "2")) {
		t.Error("expected no duplicates")
	}
	if !HasDuplicates(makeDedupEntries("A", "1", "A", "2")) {
		t.Error("expected duplicates to be detected")
	}
}
