package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func buildCastResults() []CastResult {
	return []CastResult{
		{Key: "PORT", From: " 8080 ", To: "8080", Type: CastInt, Skipped: false},
		{Key: "DEBUG", From: "1", To: "true", Type: CastBool, Skipped: false},
		{Key: "WORKERS", From: "bad", Type: CastInt, Skipped: true, Err: fmt.Errorf("cannot cast")},
	}
}

func TestCastReport_NoResults(t *testing.T) {
	rep := NewCastReport(nil)
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No cast") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestCastReport_TextWithResults(t *testing.T) {
	rep := NewCastReport(buildCastResults())
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in output")
	}
	if !strings.Contains(out, "SKIP") {
		t.Error("expected SKIP line for WORKERS")
	}
	if !strings.Contains(out, "CAST") {
		t.Error("expected CAST line")
	}
}

func TestCastReport_JSONFormat(t *testing.T) {
	rep := NewCastReport(buildCastResults())
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
	if out[2]["skipped"] != true {
		t.Error("expected skipped=true for WORKERS")
	}
}

func TestCastReport_JSONEmpty(t *testing.T) {
	rep := NewCastReport([]CastResult{})
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty JSON output")
	}
}
