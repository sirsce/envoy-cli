package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

var sampleLockEntries = func() []LockEntry {
	now := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)
	exp := now.Add(30 * time.Minute)
	return []LockEntry{
		{Key: "DB_PASSWORD", LockedBy: "alice", LockedAt: now},
		{Key: "API_KEY", LockedBy: "bob", LockedAt: now, ExpiresAt: &exp},
	}
}()

func TestLockReport_NoLocks(t *testing.T) {
	report := NewLockReport(nil)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No active locks") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestLockReport_TextWithLocks(t *testing.T) {
	report := NewLockReport(sampleLockEntries)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_PASSWORD") {
		t.Error("expected DB_PASSWORD in output")
	}
	if !strings.Contains(out, "alice") {
		t.Error("expected alice in output")
	}
	if !strings.Contains(out, "never") {
		t.Error("expected 'never' expiry for DB_PASSWORD")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Error("expected API_KEY in output")
	}
}

func TestLockReport_JSONFormat(t *testing.T) {
	report := NewLockReport(sampleLockEntries)
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestLockReport_JSONEmpty(t *testing.T) {
	report := NewLockReport([]LockEntry{})
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty array, got %d entries", len(result))
	}
}

func TestLockReport_SortedOutput(t *testing.T) {
	now := time.Now()
	entries := []LockEntry{
		{Key: "Z_KEY", LockedBy: "x", LockedAt: now},
		{Key: "A_KEY", LockedBy: "y", LockedAt: now},
	}
	report := NewLockReport(entries)
	var buf bytes.Buffer
	_ = report.WriteText(&buf)
	out := buf.String()
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Error("expected A_KEY to appear before Z_KEY in sorted output")
	}
}
