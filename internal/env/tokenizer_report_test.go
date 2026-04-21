package env

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func buildTestTokens() []Token {
	return []Token{
		{Type: TokenComment, Raw: "# header", Line: 1},
		{Type: TokenKey, Raw: "FOO=bar", Key: "FOO", Value: "bar", Line: 2},
		{Type: TokenBlank, Raw: "", Line: 3},
		{Type: TokenKey, Raw: "BAZ=qux", Key: "BAZ", Value: "qux", Line: 4},
	}
}

func TestTokenReport_Count(t *testing.T) {
	r := NewTokenReport(buildTestTokens())
	if r.Count() != 4 {
		t.Errorf("expected 4, got %d", r.Count())
	}
}

func TestTokenReport_CountByType(t *testing.T) {
	r := NewTokenReport(buildTestTokens())
	if r.CountByType(TokenKey) != 2 {
		t.Errorf("expected 2 KEY tokens")
	}
	if r.CountByType(TokenComment) != 1 {
		t.Errorf("expected 1 COMMENT token")
	}
	if r.CountByType(TokenBlank) != 1 {
		t.Errorf("expected 1 BLANK token")
	}
}

func TestTokenReport_WriteText_NoTokens(t *testing.T) {
	r := NewTokenReport(nil)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no tokens") {
		t.Errorf("expected 'no tokens', got %q", buf.String())
	}
}

func TestTokenReport_WriteText_WithTokens(t *testing.T) {
	r := NewTokenReport(buildTestTokens())
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "bar") {
		t.Errorf("expected FOO=bar in output, got %q", out)
	}
}

func TestTokenReport_JSONFormat(t *testing.T) {
	r := NewTokenReport(buildTestTokens())
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &rows); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(rows) != 4 {
		t.Errorf("expected 4 rows, got %d", len(rows))
	}
}

func TestTokenReport_JSONEmpty(t *testing.T) {
	r := NewTokenReport(nil)
	var buf bytes.Buffer
	if err := r.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got %q", buf.String())
	}
}
