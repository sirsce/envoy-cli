package env

import (
	"testing"
)

func makeCasterEntries() []Entry {
	return []Entry{
		{Key: "PORT", Value: " 8080 "},
		{Key: "RATE", Value: "3.14"},
		{Key: "DEBUG", Value: "1"},
		{Key: "NAME", Value: "  envoy  "},
		{Key: "WORKERS", Value: "not-a-number"},
	}
}

func TestCaster_CastInt(t *testing.T) {
	c := NewCaster([]CastRule{{Key: "PORT", CastTo: CastInt}})
	out, results := c.Apply(makeCasterEntries())

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].To != "8080" {
		t.Errorf("expected '8080', got %q", results[0].To)
	}
	if out[0].Value != "8080" {
		t.Errorf("entry value not updated, got %q", out[0].Value)
	}
}

func TestCaster_CastFloat(t *testing.T) {
	c := NewCaster([]CastRule{{Key: "RATE", CastTo: CastFloat}})
	_, results := c.Apply(makeCasterEntries())

	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if results[0].To != "3.14" {
		t.Errorf("expected '3.14', got %q", results[0].To)
	}
}

func TestCaster_CastBool(t *testing.T) {
	c := NewCaster([]CastRule{{Key: "DEBUG", CastTo: CastBool}})
	out, results := c.Apply(makeCasterEntries())

	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if out[2].Value != "true" {
		t.Errorf("expected 'true', got %q", out[2].Value)
	}
}

func TestCaster_CastString_TrimsSpace(t *testing.T) {
	c := NewCaster([]CastRule{{Key: "NAME", CastTo: CastString}})
	out, results := c.Apply(makeCasterEntries())

	if results[0].Err != nil {
		t.Fatalf("unexpected error: %v", results[0].Err)
	}
	if out[3].Value != "envoy" {
		t.Errorf("expected 'envoy', got %q", out[3].Value)
	}
}

func TestCaster_InvalidInt_Skipped(t *testing.T) {
	c := NewCaster([]CastRule{{Key: "WORKERS", CastTo: CastInt}})
	out, results := c.Apply(makeCasterEntries())

	if !results[0].Skipped {
		t.Error("expected result to be skipped")
	}
	if results[0].Err == nil {
		t.Error("expected non-nil error")
	}
	// original value must be preserved
	if out[4].Value != "not-a-number" {
		t.Errorf("expected original value preserved, got %q", out[4].Value)
	}
}

func TestCaster_DoesNotMutateOriginal(t *testing.T) {
	entries := makeCasterEntries()
	orig := entries[0].Value
	c := NewCaster([]CastRule{{Key: "PORT", CastTo: CastInt}})
	c.Apply(entries)

	if entries[0].Value != orig {
		t.Error("original slice was mutated")
	}
}

func TestCastedCount(t *testing.T) {
	results := []CastResult{
		{Skipped: false},
		{Skipped: true},
		{Skipped: false},
	}
	if got := CastedCount(results); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}
