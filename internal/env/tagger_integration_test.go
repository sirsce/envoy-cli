package env

import (
	"strings"
	"testing"
)

// TestTagger_WithParserAndFilter verifies that tagger works end-to-end
// with a parsed .env file to filter sensitive entries.
func TestTagger_WithParserAndFilter(t *testing.T) {
	raw := `
DB_HOST=localhost
DB_PASS=supersecret
API_KEY=abc123
APP_PORT=8080
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tg := NewTagger()
	tg.Tag("DB_PASS", "sensitive", "")
	tg.Tag("API_KEY", "sensitive", "")
	tg.Tag("APP_PORT", "optional", "")

	sensitive := tg.FilterByTag(entries, "sensitive")
	if len(sensitive) != 2 {
		t.Fatalf("expected 2 sensitive entries, got %d", len(sensitive))
	}

	keys := map[string]bool{}
	for _, e := range sensitive {
		keys[e.Key] = true
	}
	if !keys["DB_PASS"] || !keys["API_KEY"] {
		t.Error("expected DB_PASS and API_KEY in sensitive entries")
	}
}

// TestTagger_WithRedactor verifies that tagged sensitive keys are
// correctly identified and can drive redaction.
func TestTagger_WithRedactor(t *testing.T) {
	raw := `
DB_PASS=hunter2
API_KEY=tok_live_abc
APP_NAME=myapp
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tg := NewTagger()
	tg.Tag("DB_PASS", "sensitive", "")
	tg.Tag("API_KEY", "sensitive", "")

	sensitiveEntries := tg.FilterByTag(entries, "sensitive")
	sensitiveKeys := make([]string, 0, len(sensitiveEntries))
	for _, e := range sensitiveEntries {
		sensitiveKeys = append(sensitiveKeys, e.Key)
	}

	r := NewRedactor(sensitiveKeys...)
	redacted := r.RedactFull(entries)

	for _, e := range redacted {
		if tg.HasTag(e.Key, "sensitive") && e.Value != "***" {
			t.Errorf("expected %s to be redacted, got %q", e.Key, e.Value)
		}
	}

	for _, e := range redacted {
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME to be unredacted, got %q", e.Value)
		}
	}
}
