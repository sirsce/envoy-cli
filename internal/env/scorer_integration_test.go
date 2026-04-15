package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestScorer_Integration_WithParser(t *testing.T) {
	input := `DATABASE_URL=postgres://localhost/db
api_key=
DB_HOST=localhost
DB_PORT=5432
`
	p := NewParser()
	entries, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	scorer := NewScorer(IsUppercaseKey, HasEmptyValue)
	ranked := scorer.Rank(entries)

	if len(ranked) != 4 {
		t.Fatalf("expected 4 ranked entries, got %d", len(ranked))
	}
	// api_key has empty value (-1) and is lowercase (0) => score -1, should be last
	last := ranked[len(ranked)-1]
	if last.Entry.Key != "api_key" {
		t.Errorf("expected api_key to rank last, got %s", last.Entry.Key)
	}
}

func TestScorer_Integration_ReportAndExport(t *testing.T) {
	entries := []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "secret", Value: ""},
		{Key: "APP_PORT", Value: "8080"},
	}

	scorer := NewScorer(IsUppercaseKey, HasEmptyValue, HasPrefixRule("APP_"))
	ranked := scorer.Rank(entries)

	report := NewScoreReport(ranked)

	var textBuf bytes.Buffer
	if err := report.WriteText(&textBuf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := textBuf.String()
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in text report")
	}

	var jsonBuf bytes.Buffer
	if err := report.WriteJSON(&jsonBuf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(jsonBuf.String(), "APP_ENV") {
		t.Errorf("expected APP_ENV in JSON report")
	}

	// APP_ENV and APP_PORT both score 2.0 (uppercase + prefix), secret scores -1.0
	if ranked[len(ranked)-1].Entry.Key != "secret" {
		t.Errorf("expected secret to rank last, got %s", ranked[len(ranked)-1].Entry.Key)
	}
}
