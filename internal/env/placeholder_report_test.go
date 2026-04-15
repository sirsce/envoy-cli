package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildPlaceholderEntries() []Entry {
	return []Entry{
		{Key: "DSN", Value: "postgres://{{ USER }}:{{ PASS }}@localhost"},
		{Key: "HOST", Value: "localhost"},
		{Key: "URL", Value: "http://{{ HOST }}"},
	}
}

func TestPlaceholderReport_NoPlaceholders(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, false)
	entries := []Entry{{Key: "A", Value: "plain"}, {Key: "B", Value: "value"}}
	rep := NewPlaceholderReport(r, entries)
	if rep.HasPlaceholders() {
		t.Error("expected no placeholders")
	}
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "No placeholders") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPlaceholderReport_TextWithPlaceholders(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, false)
	entries := buildPlaceholderEntries()
	rep := NewPlaceholderReport(r, entries)
	if !rep.HasPlaceholders() {
		t.Fatal("expected placeholders to be detected")
	}
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DSN") {
		t.Errorf("expected DSN in output, got: %s", out)
	}
	if !strings.Contains(out, "USER") {
		t.Errorf("expected USER in output, got: %s", out)
	}
}

func TestPlaceholderReport_JSONFormat(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, false)
	entries := buildPlaceholderEntries()
	rep := NewPlaceholderReport(r, entries)
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var results []PlaceholderScanResult
	if err := json.Unmarshal(buf.Bytes(), &results); err != nil {
		t.Fatalf("JSON unmarshal error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestPlaceholderReport_JSONEmpty(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, false)
	rep := NewPlaceholderReport(r, []Entry{})
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "null") && !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array or null, got: %s", buf.String())
	}
}
