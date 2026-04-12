package env

import (
	"testing"
)

func TestSchema_Validate_AllValid(t *testing.T) {
	schema := NewSchema([]SchemaField{
		{Key: "APP_ENV", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	})
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	}
	violations := schema.Validate(entries)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestSchema_Validate_MissingRequired(t *testing.T) {
	schema := NewSchema([]SchemaField{
		{Key: "DB_URL", Required: true},
	})
	violations := schema.Validate([]Entry{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_URL" {
		t.Errorf("expected violation for DB_URL, got %q", violations[0].Key)
	}
}

func TestSchema_Validate_PatternMismatch(t *testing.T) {
	schema := NewSchema([]SchemaField{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	})
	entries := []Entry{
		{Key: "PORT", Value: "not-a-number"},
	}
	violations := schema.Validate(entries)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestSchema_Validate_RequiredWithDefault_NoViolation(t *testing.T) {
	schema := NewSchema([]SchemaField{
		{Key: "LOG_LEVEL", Required: true, Default: "info"},
	})
	violations := schema.Validate([]Entry{})
	if len(violations) != 0 {
		t.Fatalf("expected no violations when default is set, got %d", len(violations))
	}
}

func TestSchema_ApplyDefaults_FillsMissing(t *testing.T) {
	schema := NewSchema([]SchemaField{
		{Key: "LOG_LEVEL", Default: "info"},
		{Key: "TIMEOUT", Default: "30"},
	})
	entries := []Entry{
		{Key: "APP_ENV", Value: "staging"},
	}
	result := schema.ApplyDefaults(entries)

	vals := make(map[string]string)
	for _, e := range result {
		vals[e.Key] = e.Value
	}
	if vals["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", vals["LOG_LEVEL"])
	}
	if vals["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30, got %q", vals["TIMEOUT"])
	}
	if vals["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", vals["APP_ENV"])
	}
}

func TestSchema_ApplyDefaults_DoesNotOverrideExisting(t *testing.T) {
	schema := NewSchema([]SchemaField{
		{Key: "LOG_LEVEL", Default: "info"},
	})
	entries := []Entry{
		{Key: "LOG_LEVEL", Value: "debug"},
	}
	result := schema.ApplyDefaults(entries)
	for _, e := range result {
		if e.Key == "LOG_LEVEL" && e.Value != "debug" {
			t.Errorf("expected existing value to be preserved, got %q", e.Value)
		}
	}
}
