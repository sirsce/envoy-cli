package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildTestGroups() []Group {
	return []Group{
		{Name: "APP", Entries: []Entry{{Key: "APP_NAME"}, {Key: "APP_ENV"}}},
		{Name: "DB", Entries: []Entry{{Key: "DB_HOST"}}},
	}
}

func TestGroupReport_NoGroups(t *testing.T) {
	rep := NewGroupReport([]Group{})
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No groups") {
		t.Errorf("expected no-groups message, got: %s", buf.String())
	}
}

func TestGroupReport_TextWithGroups(t *testing.T) {
	rep := NewGroupReport(buildTestGroups())
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[APP]") {
		t.Errorf("expected APP group header, got: %s", out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME key, got: %s", out)
	}
	if !strings.Contains(out, "[DB] (1 keys)") {
		t.Errorf("expected DB group with count, got: %s", out)
	}
}

func TestGroupReport_JSONFormat(t *testing.T) {
	rep := NewGroupReport(buildTestGroups())
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 groups in JSON, got %d", len(result))
	}
}

func TestGroupReport_JSONEmpty(t *testing.T) {
	rep := NewGroupReport([]Group{})
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty JSON array")
	}
}
