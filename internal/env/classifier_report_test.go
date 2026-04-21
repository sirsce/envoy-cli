package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildTestClassifyResults() []ClassifyResult {
	return []ClassifyResult{
		{Key: "API_KEY", Label: "secret", Confidence: 0.9},
		{Key: "DB_HOST", Label: "database", Confidence: 0.9},
		{Key: "APP_NAME", Label: "generic", Confidence: 0.5},
	}
}

func TestClassifierReport_NoResults(t *testing.T) {
	r := NewClassifierReport(nil)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestClassifierReport_TextWithResults(t *testing.T) {
	r := NewClassifierReport(buildTestClassifyResults())
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output")
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("expected secret label in output")
	}
}

func TestClassifierReport_JSONFormat(t *testing.T) {
	r := NewClassifierReport(buildTestClassifyResults())
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []ClassifyResult
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 results, got %d", len(out))
	}
}

func TestClassifierReport_Summary(t *testing.T) {
	r := NewClassifierReport(buildTestClassifyResults())
	summary := r.Summary()
	if summary["secret"] != 1 {
		t.Errorf("expected 1 secret, got %d", summary["secret"])
	}
	if summary["generic"] != 1 {
		t.Errorf("expected 1 generic, got %d", summary["generic"])
	}
}

func TestClassifierReport_JSONEmpty(t *testing.T) {
	r := NewClassifierReport([]ClassifyResult{})
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[") {
		t.Errorf("expected JSON array in output")
	}
}
