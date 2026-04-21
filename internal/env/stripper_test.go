package env

import (
	"bytes"
	"strings"
	"testing"
)

func makeStripperEntries() []Entry {
	return []Entry{
		{Key: "API_KEY", Value: "secret_abc123"},
		{Key: "DB_PASS", Value: "hunter2"},
		{Key: "HOST", Value: "localhost"},
		{Key: "TOKEN", Value: "bearer_xyz"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestStripper_NoPatterns(t *testing.T) {
	s := NewStripper()
	out, results := s.Strip(makeStripperEntries())
	if len(out) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(out))
	}
	if StrippedCount(results) != 0 {
		t.Fatalf("expected 0 stripped, got %d", StrippedCount(results))
	}
}

func TestStripper_ByValuePrefix(t *testing.T) {
	s := NewStripper(WithStripValuePrefix("secret_"), WithStripValuePrefix("bearer_"))
	out, results := s.Strip(makeStripperEntries())
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if StrippedCount(results) != 2 {
		t.Fatalf("expected 2 stripped, got %d", StrippedCount(results))
	}
	for _, e := range out {
		if e.Key == "API_KEY" || e.Key == "TOKEN" {
			t.Errorf("stripped key %q should not be in output", e.Key)
		}
	}
}

func TestStripper_ByValueContains(t *testing.T) {
	s := NewStripper(WithStripValueContains("hunter"))
	out, results := s.Strip(makeStripperEntries())
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	if StrippedCount(results) != 1 {
		t.Fatalf("expected 1 stripped, got %d", StrippedCount(results))
	}
}

func TestStripper_BlankInstead(t *testing.T) {
	s := NewStripper(WithStripValuePrefix("secret_"), WithBlankInstead())
	out, results := s.Strip(makeStripperEntries())
	if len(out) != 5 {
		t.Fatalf("expected 5 entries (blanked, not removed), got %d", len(out))
	}
	for _, e := range out {
		if e.Key == "API_KEY" && e.Value != "" {
			t.Errorf("expected API_KEY to be blanked")
		}
	}
	if StrippedCount(results) != 1 {
		t.Fatalf("expected 1 in results, got %d", StrippedCount(results))
	}
}

func TestStripper_ReportText(t *testing.T) {
	s := NewStripper(WithStripValuePrefix("secret_"))
	_, results := s.Strip(makeStripperEntries())
	report := NewStripReport(results)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "API_KEY") {
		t.Errorf("expected API_KEY in text report, got: %s", buf.String())
	}
}

func TestStripper_ReportJSON(t *testing.T) {
	s := NewStripper(WithStripValueContains("hunter"))
	_, results := s.Strip(makeStripperEntries())
	report := NewStripReport(results)
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "DB_PASS") {
		t.Errorf("expected DB_PASS in JSON report, got: %s", buf.String())
	}
}
