package env

import (
	"testing"
)

func makeRenamerEntries(pairs ...string) []Entry {
	entries := make([]Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestRenamer_AppliesRule(t *testing.T) {
	entries := makeRenamerEntries("OLD_KEY", "hello")
	r := NewRenamer([]RenameRule{{From: "OLD_KEY", To: "NEW_KEY"}})

	out, results, err := r.Apply(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Errorf("expected rule to be applied")
	}
	if out[0].Key != "NEW_KEY" {
		t.Errorf("expected key NEW_KEY, got %s", out[0].Key)
	}
	if out[0].Value != "hello" {
		t.Errorf("expected value hello, got %s", out[0].Value)
	}
}

func TestRenamer_SourceNotFound(t *testing.T) {
	entries := makeRenamerEntries("SOME_KEY", "val")
	r := NewRenamer([]RenameRule{{From: "MISSING", To: "TARGET"}})

	_, results, err := r.Apply(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Errorf("expected rule NOT to be applied")
	}
	if results[0].Reason != "source key not found" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestRenamer_DestinationExists(t *testing.T) {
	entries := makeRenamerEntries("A", "1", "B", "2")
	r := NewRenamer([]RenameRule{{From: "A", To: "B"}})

	_, results, err := r.Apply(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Errorf("expected rule NOT to be applied")
	}
	if results[0].Reason != "destination key already exists" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestRenamer_EmptyRuleReturnsError(t *testing.T) {
	entries := makeRenamerEntries("KEY", "val")
	r := NewRenamer([]RenameRule{{From: "", To: "NEW"}})

	_, _, err := r.Apply(entries)
	if err == nil {
		t.Error("expected error for empty From field")
	}
}

func TestRenamer_AppliedCount(t *testing.T) {
	results := []RenameResult{
		{Applied: true},
		{Applied: false},
		{Applied: true},
	}
	if got := AppliedCount(results); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestRenamer_DoesNotMutateOnSkip(t *testing.T) {
	entries := makeRenamerEntries("X", "original")
	r := NewRenamer([]RenameRule{{From: "MISSING", To: "Y"}})

	out, _, _ := r.Apply(entries)
	if out[0].Key != "X" {
		t.Errorf("expected key X to be unchanged, got %s", out[0].Key)
	}
}
