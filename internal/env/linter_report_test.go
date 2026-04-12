package env

import (
	"strings"
	"testing"
)

func TestLintReport_TextNoViolations(t *testing.T) {
	r := NewLintReport(ReportFormatText)
	var sb strings.Builder
	if err := r.Write(&sb, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No lint violations") {
		t.Errorf("expected no-violations message, got: %q", sb.String())
	}
}

func TestLintReport_TextWithViolations(t *testing.T) {
	r := NewLintReport(ReportFormatText)
	results := []LintResult{
		{Entry: Entry{Key: "secret", Value: ""}, Rule: "no-empty-value", Message: "value should not be empty: key=\"secret\""},
	}
	var sb strings.Builder
	if err := r.Write(&sb, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "no-empty-value") {
		t.Errorf("expected rule name in output, got: %q", out)
	}
	if !strings.Contains(out, "secret") {
		t.Errorf("expected key in output, got: %q", out)
	}
}

func TestLintReport_JSONFormat(t *testing.T) {
	r := NewLintReport(ReportFormatJSON)
	results := []LintResult{
		{Entry: Entry{Key: "MY_KEY", Value: ""}, Rule: "no-empty-value", Message: "value should not be empty: key=\"MY_KEY\""},
	}
	var sb strings.Builder
	if err := r.Write(&sb, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "\"rule\"") {
		t.Errorf("expected JSON key 'rule', got: %q", out)
	}
	if !strings.Contains(out, "no-empty-value") {
		t.Errorf("expected rule value in JSON, got: %q", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("expected JSON array, got: %q", out)
	}
}

func TestLintReport_JSONEmpty(t *testing.T) {
	r := NewLintReport(ReportFormatJSON)
	var sb strings.Builder
	if err := r.Write(&sb, []LintResult{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(sb.String())
	if out != "[]" {
		t.Errorf("expected empty JSON array '[]', got: %q", out)
	}
}

func TestLintReport_TextEmptySlice(t *testing.T) {
	// Passing an empty (non-nil) slice should behave the same as nil: no violations.
	r := NewLintReport(ReportFormatText)
	var sb strings.Builder
	if err := r.Write(&sb, []LintResult{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No lint violations") {
		t.Errorf("expected no-violations message for empty slice, got: %q", sb.String())
	}
}
