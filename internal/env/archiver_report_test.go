package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func buildTestArchiveEntries() []ArchiveEntry {
	t0 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []ArchiveEntry{
		{Version: 1, Label: "baseline", Entries: []Entry{{Key: "A"}, {Key: "B"}}, ArchivedAt: t0},
		{Version: 2, Label: "production", Entries: []Entry{{Key: "X"}}, ArchivedAt: t0.Add(24 * time.Hour)},
	}
}

func TestArchiveReport_NoEntries(t *testing.T) {
	r := NewArchiveReport(nil)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No archives") {
		t.Errorf("expected empty message, got: %q", buf.String())
	}
}

func TestArchiveReport_TextWithEntries(t *testing.T) {
	r := NewArchiveReport(buildTestArchiveEntries())
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "baseline") {
		t.Errorf("expected 'baseline' in output, got: %q", out)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected 'production' in output, got: %q", out)
	}
	if !strings.Contains(out, "2 keys") {
		t.Errorf("expected '2 keys' in output, got: %q", out)
	}
}

func TestArchiveReport_JSONFormat(t *testing.T) {
	r := NewArchiveReport(buildTestArchiveEntries())
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []archiveReportJSON
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Label != "baseline" || out[0].KeyCount != 2 {
		t.Errorf("unexpected first entry: %+v", out[0])
	}
	if out[1].Label != "production" || out[1].KeyCount != 1 {
		t.Errorf("unexpected second entry: %+v", out[1])
	}
}

func TestArchiveReport_JSONEmpty(t *testing.T) {
	r := NewArchiveReport([]ArchiveEntry{})
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []archiveReportJSON
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty array, got %d entries", len(out))
	}
}
