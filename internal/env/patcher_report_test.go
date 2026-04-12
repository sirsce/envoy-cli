package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func samplePatchResults() []PatchResult {
	return []PatchResult{
		{Instruction: PatchInstruction{Op: PatchSet, Key: "HOST", Value: "db"}, Applied: true, Note: "updated"},
		{Instruction: PatchInstruction{Op: PatchDelete, Key: "GHOST"}, Applied: false, Note: "key not found"},
		{Instruction: PatchInstruction{Op: PatchRename, Key: "PORT", NewKey: "DB_PORT"}, Applied: true, Note: "renamed"},
	}
}

func TestPatchReport_Applied(t *testing.T) {
	r := NewPatchReport(samplePatchResults())
	if len(r.Applied()) != 2 {
		t.Errorf("expected 2 applied, got %d", len(r.Applied()))
	}
}

func TestPatchReport_Skipped(t *testing.T) {
	r := NewPatchReport(samplePatchResults())
	if len(r.Skipped()) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(r.Skipped()))
	}
}

func TestPatchReport_WriteText_WithResults(t *testing.T) {
	r := NewPatchReport(samplePatchResults())
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Error("expected [OK] in output")
	}
	if !strings.Contains(out, "[SKIP]") {
		t.Error("expected [SKIP] in output")
	}
	if !strings.Contains(out, "HOST") {
		t.Error("expected HOST in output")
	}
}

func TestPatchReport_WriteText_Empty(t *testing.T) {
	r := NewPatchReport(nil)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "No patch") {
		t.Error("expected empty message")
	}
}

func TestPatchReport_WriteJSON(t *testing.T) {
	r := NewPatchReport(samplePatchResults())
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(rows))
	}
	if rows[1]["applied"].(bool) {
		t.Error("expected second result applied=false")
	}
}

func TestPatchReport_JSONEmpty(t *testing.T) {
	r := NewPatchReport([]PatchResult{})
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 0 {
		t.Errorf("expected empty array, got %d items", len(rows))
	}
}
