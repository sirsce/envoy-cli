package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildTestLabeler() *Labeler {
	l := NewLabeler()
	_ = l.Set("DB_HOST", "env", "production")
	_ = l.Set("DB_HOST", "owner", "backend")
	_ = l.Set("API_KEY", "env", "production")
	return l
}

func TestLabelReport_NoLabels(t *testing.T) {
	l := NewLabeler()
	r := NewLabelReport(l)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No labels defined.") {
		t.Errorf("expected empty message, got: %q", buf.String())
	}
}

func TestLabelReport_TextWithLabels(t *testing.T) {
	l := buildTestLabeler()
	r := NewLabelReport(l)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST:") {
		t.Errorf("expected DB_HOST in output, got: %q", out)
	}
	if !strings.Contains(out, "env=production") {
		t.Errorf("expected env=production in output, got: %q", out)
	}
	if !strings.Contains(out, "owner=backend") {
		t.Errorf("expected owner=backend in output, got: %q", out)
	}
}

func TestLabelReport_JSONFormat(t *testing.T) {
	l := buildTestLabeler()
	r := NewLabelReport(l)
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []LabelReportEntry
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestLabelReport_JSONEmpty(t *testing.T) {
	l := NewLabeler()
	r := NewLabelReport(l)
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []LabelReportEntry
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty array, got %d entries", len(result))
	}
}

func TestLabelReport_SortedOutput(t *testing.T) {
	l := NewLabeler()
	_ = l.Set("Z_KEY", "env", "prod")
	_ = l.Set("A_KEY", "env", "prod")
	r := NewLabelReport(l)
	var buf bytes.Buffer
	_ = r.WriteText(&buf)
	out := buf.String()
	azIdx := strings.Index(out, "A_KEY")
	zzIdx := strings.Index(out, "Z_KEY")
	if azIdx > zzIdx {
		t.Error("expected A_KEY to appear before Z_KEY in sorted output")
	}
}
