package env

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func sampleRotationRecords() []RotationRecord {
	return []RotationRecord{
		{Key: "DB_PASSWORD", OldValue: "old1", NewValue: "new1", RotatedAt: time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)},
		{Key: "API_SECRET", OldValue: "old2", NewValue: "new2", RotatedAt: time.Date(2024, 3, 1, 10, 0, 1, 0, time.UTC)},
	}
}

func TestRotationReport_NoRotations(t *testing.T) {
	rr := NewRotationReport(nil)
	if rr.HasRotations() {
		t.Error("expected HasRotations to be false")
	}
	var buf bytes.Buffer
	if err := rr.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "No keys rotated") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestRotationReport_TextWithRotations(t *testing.T) {
	rr := NewRotationReport(sampleRotationRecords())
	if !rr.HasRotations() {
		t.Error("expected HasRotations to be true")
	}
	var buf bytes.Buffer
	if err := rr.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Rotated 2 key(s)") {
		t.Errorf("missing header in output: %q", out)
	}
	if !strings.Contains(out, "DB_PASSWORD") || !strings.Contains(out, "API_SECRET") {
		t.Errorf("missing key names in output: %q", out)
	}
}

func TestRotationReport_JSONFormat(t *testing.T) {
	rr := NewRotationReport(sampleRotationRecords())
	var buf bytes.Buffer
	if err := rr.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"key"`) {
		t.Errorf("expected JSON key field, got: %q", out)
	}
	if !strings.Contains(out, `"rotated_at"`) {
		t.Errorf("expected JSON rotated_at field, got: %q", out)
	}
	if strings.Contains(out, "old1") || strings.Contains(out, "new1") {
		t.Error("old/new values must not appear in JSON report")
	}
}

func TestRotationReport_JSONEmpty(t *testing.T) {
	rr := NewRotationReport([]RotationRecord{})
	var buf bytes.Buffer
	if err := rr.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got: %q", buf.String())
	}
}
