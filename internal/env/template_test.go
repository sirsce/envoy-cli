package env

import (
	"strings"
	"testing"
)

func TestNewTemplate_RequiredAndDefaults(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: ""},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "development"},
	}
	tmpl := NewTemplate(entries)

	if !tmpl.required["DB_HOST"] {
		t.Error("expected DB_HOST to be required")
	}
	if tmpl.required["DB_PORT"] {
		t.Error("expected DB_PORT to not be required")
	}
	if tmpl.defaults["DB_PORT"] != "5432" {
		t.Errorf("expected default DB_PORT=5432, got %s", tmpl.defaults["DB_PORT"])
	}
}

func TestTemplate_Apply_Success(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: ""},
		{Key: "DB_PORT", Value: "5432"},
	}
	tmpl := NewTemplate(entries)

	result, err := tmpl.Apply(map[string]string{"DB_HOST": "localhost"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Value != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", result[0].Value)
	}
	if result[1].Value != "5432" {
		t.Errorf("expected DB_PORT=5432 (default), got %s", result[1].Value)
	}
}

func TestTemplate_Apply_MissingRequired(t *testing.T) {
	entries := []Entry{
		{Key: "SECRET_KEY", Value: ""},
		{Key: "API_URL", Value: ""},
	}
	tmpl := NewTemplate(entries)

	_, err := tmpl.Apply(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing required keys")
	}
	if !strings.Contains(err.Error(), "SECRET_KEY") {
		t.Errorf("expected error to mention SECRET_KEY, got: %v", err)
	}
}

func TestTemplate_Render(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "", Comment: "Database hostname"},
		{Key: "DB_PORT", Value: "5432"},
	}
	tmpl := NewTemplate(entries)
	out := tmpl.Render()

	if !strings.Contains(out, "# Database hostname") {
		t.Error("expected comment in render output")
	}
	if !strings.Contains(out, "DB_HOST=") {
		t.Error("expected DB_HOST= in render output")
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Error("expected DB_PORT=5432 in render output")
	}
}

func TestTemplate_RequiredKeys(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: ""},
		{Key: "B", Value: "default"},
		{Key: "C", Value: ""},
	}
	tmpl := NewTemplate(entries)
	keys := tmpl.RequiredKeys()

	if len(keys) != 2 {
		t.Errorf("expected 2 required keys, got %d", len(keys))
	}
}
