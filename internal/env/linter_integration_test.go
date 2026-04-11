package env

import (
	"strings"
	"testing"
)

func TestLinter_WithParser(t *testing.T) {
	raw := `DATABASE_URL=postgres://localhost/mydb
APP_PORT=3000
SECRET=
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	l := NewLinter()
	results := l.Lint(entries)

	if len(results) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(results), results)
	}
	if results[0].Rule != "no-empty-value" {
		t.Errorf("expected 'no-empty-value', got %q", results[0].Rule)
	}
}

func TestLinter_MixedViolations(t *testing.T) {
	raw := `lowercase_key=value
GOOD_KEY=
ANOTHER=fine
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	l := NewLinter()
	results := l.Lint(entries)

	ruleHits := map[string]int{}
	for _, r := range results {
		ruleHits[r.Rule]++
	}

	if ruleHits["uppercase-key"] != 1 {
		t.Errorf("expected 1 uppercase-key violation, got %d", ruleHits["uppercase-key"])
	}
	if ruleHits["no-empty-value"] != 1 {
		t.Errorf("expected 1 no-empty-value violation, got %d", ruleHits["no-empty-value"])
	}
}

func TestLinter_AllClean(t *testing.T) {
	raw := `APP_NAME=envoy
APP_ENV=production
PORT=443
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	l := NewLinter()
	results := l.Lint(entries)

	if HasViolations(results) {
		t.Errorf("expected no violations, got: %+v", results)
	}
}
