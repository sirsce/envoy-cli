package env

import (
	"testing"
)

func TestSorter_Alpha(t *testing.T) {
	entries := []Entry{
		{Key: "ZEBRA", Value: "1"},
		{Key: "APPLE", Value: "2"},
		{Key: "MANGO", Value: "3"},
	}

	s := NewSorter(SortAlpha)
	result := s.Sort(entries)

	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("index %d: expected key %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestSorter_AlphaDesc(t *testing.T) {
	entries := []Entry{
		{Key: "APPLE", Value: "1"},
		{Key: "ZEBRA", Value: "2"},
		{Key: "MANGO", Value: "3"},
	}

	s := NewSorter(SortAlphaDesc)
	result := s.Sort(entries)

	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("index %d: expected key %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestSorter_ByGroup(t *testing.T) {
	entries := []Entry{
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "APP_ENV", Value: "prod"},
		{Key: "LOG_LEVEL", Value: "info"},
	}

	s := NewSorter(SortByGroup)
	result := s.Sort(entries)

	expected := []string{"APP_ENV", "APP_NAME", "DB_HOST", "DB_PORT", "LOG_LEVEL"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("index %d: expected key %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestSorter_DoesNotMutateOriginal(t *testing.T) {
	entries := []Entry{
		{Key: "Z", Value: "1"},
		{Key: "A", Value: "2"},
	}

	s := NewSorter(SortAlpha)
	_ = s.Sort(entries)

	if entries[0].Key != "Z" {
		t.Error("original slice was mutated")
	}
}

func TestSorter_Empty(t *testing.T) {
	s := NewSorter(SortAlpha)
	result := s.Sort([]Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
