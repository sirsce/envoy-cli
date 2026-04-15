package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildTestScoredEntries() []ScoredEntry {
	return []ScoredEntry{
		{Entry: Entry{Key: "DATABASE_URL", Value: "postgres://localhost"}, Score: 2.0},
		{Entry: Entry{Key: "api_key", Value: ""}, Score: -1.0},
	}
}

func TestScoreReport_NoEntries(t *testing.T) {
	report := NewScoreReport(nil)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No scored entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestScoreReport_TextWithEntries(t *testing.T) {
	report := NewScoreReport(buildTestScoredEntries())
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DATABASE_URL") {
		t.Errorf("expected DATABASE_URL in output")
	}
	if !strings.Contains(out, "2.00") {
		t.Errorf("expected score 2.00 in output")
	}
	if !strings.Contains(out, "-1.00") {
		t.Errorf("expected score -1.00 in output")
	}
}

func TestScoreReport_JSONFormat(t *testing.T) {
	report := NewScoreReport(buildTestScoredEntries())
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []struct {
		Key   string  `json:"key"`
		Value string  `json:"value"`
		Score float64 `json:"score"`
	}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].Key != "DATABASE_URL" || rows[0].Score != 2.0 {
		t.Errorf("unexpected first row: %+v", rows[0])
	}
}

func TestScoreReport_JSONEmpty(t *testing.T) {
	report := NewScoreReport([]ScoredEntry{})
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected empty array, got %d items", len(rows))
	}
}
