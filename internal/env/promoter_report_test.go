package env_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/yourusername/envoy-cli/internal/env"
)

func buildTestPromoteReport() *env.PromoteReport {
	results := []env.PromoteResult{
		{Key: "APP_ENV", FromScope: "development", ToScope: "staging", Applied: true, Overwritten: false},
		{Key: "DB_URL", FromScope: "development", ToScope: "staging", Applied: true, Overwritten: true},
		{Key: "SECRET_KEY", FromScope: "development", ToScope: "staging", Applied: false, Overwritten: false, Reason: "already exists"},
	}
	return env.NewPromoteReport(results)
}

func TestPromoteReport_NoResults(t *testing.T) {
	report := env.NewPromoteReport(nil)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if output != "No promotions performed.\n" {
		t.Errorf("expected no-promotions message, got: %q", output)
	}
}

func TestPromoteReport_TextWithResults(t *testing.T) {
	report := buildTestPromoteReport()
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()

	for _, key := range []string{"APP_ENV", "DB_URL", "SECRET_KEY"} {
		if !bytes.Contains(buf.Bytes(), []byte(key)) {
			t.Errorf("expected key %q in text output, got:\n%s", key, output)
		}
	}
	if !bytes.Contains(buf.Bytes(), []byte("overwritten")) {
		t.Errorf("expected 'overwritten' label in output, got:\n%s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("skipped")) {
		t.Errorf("expected 'skipped' label in output, got:\n%s", output)
	}
}

func TestPromoteReport_JSONFormat(t *testing.T) {
	report := buildTestPromoteReport()
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results in JSON, got %d", len(results))
	}

	first := results[0]
	if first["key"] != "APP_ENV" {
		t.Errorf("expected first key to be APP_ENV, got %v", first["key"])
	}
	if first["applied"] != true {
		t.Errorf("expected first result to be applied")
	}
}

func TestPromoteReport_JSONEmpty(t *testing.T) {
	report := env.NewPromoteReport([]env.PromoteResult{})
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected empty JSON array, got %d items", len(results))
	}
}

func TestPromoteReport_Summary(t *testing.T) {
	report := buildTestPromoteReport()
	applied, skipped := report.Summary()
	if applied != 2 {
		t.Errorf("expected 2 applied, got %d", applied)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", skipped)
	}
}
