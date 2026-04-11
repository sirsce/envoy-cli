package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func sampleEvents() []AuditEvent {
	ts := fixedTime()
	return []AuditEvent{
		{Timestamp: ts, Action: AuditActionAdd, Key: "NEW_KEY", NewValue: "hello"},
		{Timestamp: ts, Action: AuditActionRemove, Key: "OLD_KEY", OldValue: "bye"},
		{Timestamp: ts, Action: AuditActionChange, Key: "CHANGED", OldValue: "a", NewValue: "b"},
	}
}

func TestAuditLog_HasEvents(t *testing.T) {
	log := NewAuditLog(nil)
	if log.HasEvents() {
		t.Error("expected no events")
	}
	log2 := NewAuditLog(sampleEvents())
	if !log2.HasEvents() {
		t.Error("expected events")
	}
}

func TestAuditLog_FilterByAction(t *testing.T) {
	log := NewAuditLog(sampleEvents())
	added := log.FilterByAction(AuditActionAdd)
	if len(added) != 1 || added[0].Key != "NEW_KEY" {
		t.Errorf("unexpected filter result: %+v", added)
	}
}

func TestAuditLog_WriteText_NoEvents(t *testing.T) {
	var buf bytes.Buffer
	log := NewAuditLog(nil)
	if err := log.WriteText(&buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestAuditLog_WriteText_WithEvents(t *testing.T) {
	var buf bytes.Buffer
	log := NewAuditLog(sampleEvents())
	if err := log.WriteText(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "ADD NEW_KEY") {
		t.Errorf("expected ADD line, got: %q", out)
	}
	if !strings.Contains(out, "REMOVE OLD_KEY") {
		t.Errorf("expected REMOVE line, got: %q", out)
	}
}

func TestAuditLog_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	log := NewAuditLog(sampleEvents())
	if err := log.WriteJSON(&buf); err != nil {
		t.Fatal(err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 JSON entries, got %d", len(result))
	}
	if result[0]["action"] != "add" {
		t.Errorf("expected action=add, got %v", result[0]["action"])
	}
}

func TestAuditLog_WriteJSON_Empty(t *testing.T) {
	var buf bytes.Buffer
	log := NewAuditLog(nil)
	if err := log.WriteJSON(&buf); err != nil {
		t.Fatal(err)
	}
	_ = time.Now() // ensure time import used
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty JSON output")
	}
}
