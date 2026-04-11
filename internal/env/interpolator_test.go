package env

import (
	"testing"
)

func TestInterpolator_NoReferences(t *testing.T) {
	i := NewInterpolator()
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
	}
	out, err := i.Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" || out[1].Value != "5432" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestInterpolator_BraceStyle(t *testing.T) {
	i := NewInterpolator()
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "DSN", Value: "postgres://${HOST}:5432/db"},
	}
	out, err := i.Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "postgres://localhost:5432/db" {
		t.Errorf("got %q", out[1].Value)
	}
}

func TestInterpolator_DollarStyle(t *testing.T) {
	i := NewInterpolator()
	entries := []Entry{
		{Key: "USER", Value: "admin"},
		{Key: "GREETING", Value: "hello $USER"},
	}
	out, err := i.Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "hello admin" {
		t.Errorf("got %q", out[1].Value)
	}
}

func TestInterpolator_UndefinedVariable_Error(t *testing.T) {
	i := NewInterpolator()
	entries := []Entry{
		{Key: "DSN", Value: "postgres://${MISSING_HOST}:5432"},
	}
	_, err := i.Interpolate(entries)
	if err == nil {
		t.Fatal("expected error for undefined variable")
	}
}

func TestInterpolator_AllowMissing_LeavesReference(t *testing.T) {
	i := NewInterpolator(WithAllowMissing())
	entries := []Entry{
		{Key: "URL", Value: "http://${UNDEFINED}/path"},
	}
	out, err := i.Interpolate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "http://${UNDEFINED}/path" {
		t.Errorf("expected reference preserved, got %q", out[0].Value)
	}
}

func TestInterpolator_DoesNotMutateInput(t *testing.T) {
	i := NewInterpolator()
	original := []Entry{
		{Key: "BASE", Value: "example.com"},
		{Key: "URL", Value: "https://${BASE}"},
	}
	_, err := i.Interpolate(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if original[1].Value != "https://${BASE}" {
		t.Errorf("input was mutated: got %q", original[1].Value)
	}
}
