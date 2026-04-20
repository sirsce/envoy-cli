package env

import (
	"testing"
)

func makeSplitterEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "dev"},
		{Key: "SECRET", Value: "abc"},
	}
}

func TestSplitter_ByPrefix(t *testing.T) {
	s := NewSplitter(SplitByPrefix, 0)
	results, err := s.Split(makeSplitterEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bucketNames := map[string]bool{}
	for _, r := range results {
		bucketNames[r.Name] = true
	}
	if !bucketNames["DB"] {
		t.Error("expected bucket DB")
	}
	if !bucketNames["APP"] {
		t.Error("expected bucket APP")
	}
	if !bucketNames["_"] && !bucketNames["SECRET"] {
		// SECRET has no underscore, groupPrefix returns ""
		t.Error("expected fallback bucket for SECRET")
	}
}

func TestSplitter_ByPrefix_Counts(t *testing.T) {
	s := NewSplitter(SplitByPrefix, 0)
	results, _ := s.Split(makeSplitterEntries())
	for _, r := range results {
		if r.Name == "DB" && len(r.Entries) != 2 {
			t.Errorf("DB bucket: want 2 entries, got %d", len(r.Entries))
		}
		if r.Name == "APP" && len(r.Entries) != 2 {
			t.Errorf("APP bucket: want 2 entries, got %d", len(r.Entries))
		}
	}
}

func TestSplitter_ByFirstChar(t *testing.T) {
	s := NewSplitter(SplitByFirstChar, 0)
	entries := []Entry{
		{Key: "ALPHA", Value: "1"},
		{Key: "BETA", Value: "2"},
		{Key: "BRAVO", Value: "3"},
	}
	results, err := s.Split(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("want 2 buckets (A, B), got %d", len(results))
	}
}

func TestSplitter_Even(t *testing.T) {
	s := NewSplitter(SplitEven, 2)
	results, err := s.Split(makeSplitterEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("want 2 buckets, got %d", len(results))
	}
	total := len(results[0].Entries) + len(results[1].Entries)
	if total != 5 {
		t.Errorf("want 5 total entries, got %d", total)
	}
}

func TestSplitter_Even_MinBuckets(t *testing.T) {
	s := NewSplitter(SplitEven, 0) // should clamp to 1
	results, err := s.Split(makeSplitterEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("want 1 bucket, got %d", len(results))
	}
	if len(results[0].Entries) != 5 {
		t.Errorf("want 5 entries in single bucket, got %d", len(results[0].Entries))
	}
}

func TestSplitter_Empty(t *testing.T) {
	s := NewSplitter(SplitByPrefix, 0)
	results, err := s.Split([]Entry{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("want 0 results for empty input, got %d", len(results))
	}
}
