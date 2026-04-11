package env

import (
	"strings"
	"testing"
)

// TestTransformer_WithParser verifies transformer works on parser output.
func TestTransformer_WithParser(t *testing.T) {
	input := `
DB_HOST=  localhost  
DB_PORT=  5432  
APP_ENV=  production  
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tr := NewTransformer(TrimSpaceTransform())
	result := tr.Apply(entries)

	m := make(map[string]string)
	for _, e := range result {
		m[e.Key] = e.Value
	}

	if m["DB_HOST"] != "localhost" {
		t.Errorf("expected 'localhost', got %q", m["DB_HOST"])
	}
	if m["DB_PORT"] != "5432" {
		t.Errorf("expected '5432', got %q", m["DB_PORT"])
	}
	if m["APP_ENV"] != "production" {
		t.Errorf("expected 'production', got %q", m["APP_ENV"])
	}
}

// TestTransformer_WithFilterAndExport verifies transformer composes with filter and exporter.
func TestTransformer_WithFilterAndExport(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "  localhost  "},
		{Key: "DB_PORT", Value: "  5432  "},
		{Key: "APP_SECRET", Value: "  s3cr3t  "},
	}

	filtered := FilterByPrefix(entries, "DB_")
	tr := NewTransformer(TrimSpaceTransform(), ToUpperTransform())
	transformed := tr.Apply(filtered)

	if len(transformed) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(transformed))
	}
	for _, e := range transformed {
		if strings.Contains(e.Value, " ") {
			t.Errorf("value %q should have no spaces after trim", e.Value)
		}
		if e.Value != strings.ToUpper(e.Value) {
			t.Errorf("value %q should be uppercase", e.Value)
		}
	}
}
