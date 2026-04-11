package env

import (
	"testing"
)

func TestLinter_NoViolations(t *testing.T) {
	l := NewLinter()
	entries := []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "APP_PORT", Value: "8080"},
	}
	results := l.Lint(entries)
	if len(results) != 0 {
		t.Errorf("expected no violations, got %d", len(results))
	}
}

func TestLinter_EmptyValue(t *testing.T) {
	l := NewLinter()
	entries := []Entry{
		{Key: "SECRET_KEY", Value: ""},
	}
	results := l.Lint(entries)
	if !HasViolations(results) {
		t.Fatal("expected violation for empty value")
	}
	if results[0].Rule != "no-empty-value" {
		t.Errorf("expected rule 'no-empty-value', got %q", results[0].Rule)
	}
}

func TestLinter_LowercaseKey(t *testing.T) {
	l := NewLinter()
	entries := []Entry{
		{Key: "database_url", Value: "something"},
	}
	results := l.Lint(entries)
	var found bool
	for _, r := range results {
		if r.Rule == "uppercase-key" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'uppercase-key' violation")
	}
}

func TestLinter_WhitespaceInKey(t *testing.T) {
	l := NewLinter()
	entries := []Entry{
		{Key: "MY KEY", Value: "value"},
	}
	results := l.Lint(entries)
	var found bool
	for _, r := range results {
		if r.Rule == "no-whitespace-in-key" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'no-whitespace-in-key' violation")
	}
}

func TestLinter_CustomRule(t *testing.T) {
	l := NewLinter()
	l.AddRule(LintRule{
		Name:    "no-localhost",
		Message: "value should not reference localhost",
		Check: func(e Entry) bool {
			return !contains(e.Value, "localhost")
		},
	})
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost:5432"},
	}
	results := l.Lint(entries)
	var found bool
	for _, r := range results {
		if r.Rule == "no-localhost" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'no-localhost' violation")
	}
}

func TestHasViolations_Empty(t *testing.T) {
	if HasViolations(nil) {
		t.Error("expected false for nil results")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
