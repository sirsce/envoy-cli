package env

import (
	"testing"
)

func testEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "db"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
		{Key: "APP_SECRET", Value: "topsecret"},
	}
}

func TestFilter_Prefix(t *testing.T) {
	f := NewFilter(FilterOptions{Prefix: "APP_"})
	got := f.Apply(testEntries())
	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key[:4] != "APP_" {
			t.Errorf("unexpected key %q", e.Key)
		}
	}
}

func TestFilter_Suffix(t *testing.T) {
	f := NewFilter(FilterOptions{Suffix: "_PORT"})
	got := f.Apply(testEntries())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilter_KeySubstr(t *testing.T) {
	f := NewFilter(FilterOptions{KeySubstr: "SECRET"})
	got := f.Apply(testEntries())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilter_ExcludeKeys(t *testing.T) {
	f := NewFilter(FilterOptions{ExcludeKeys: []string{"DB_HOST", "DB_PORT"}})
	got := f.Apply(testEntries())
	if len(got) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(got))
	}
	for _, e := range got {
		if e.Key == "DB_HOST" || e.Key == "DB_PORT" {
			t.Errorf("excluded key %q still present", e.Key)
		}
	}
}

func TestFilter_Combined(t *testing.T) {
	f := NewFilter(FilterOptions{Prefix: "APP_", ExcludeKeys: []string{"APP_SECRET"}})
	got := f.Apply(testEntries())
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestFilterByPrefix_Convenience(t *testing.T) {
	got := FilterByPrefix(testEntries(), "DB_")
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestStripPrefix(t *testing.T) {
	entries := []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "OTHER", Value: "val"},
	}
	got := StripPrefix(entries, "APP_")
	if got[0].Key != "HOST" {
		t.Errorf("expected HOST, got %q", got[0].Key)
	}
	if got[1].Key != "PORT" {
		t.Errorf("expected PORT, got %q", got[1].Key)
	}
	if got[2].Key != "OTHER" {
		t.Errorf("expected OTHER unchanged, got %q", got[2].Key)
	}
}
