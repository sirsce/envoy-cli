package env

import (
	"testing"
)

func TestMasker_IsSensitive(t *testing.T) {
	m := NewMasker(MaskFull, "API_KEY", "DB_PASSWORD")

	if !m.IsSensitive("API_KEY") {
		t.Error("expected API_KEY to be sensitive")
	}
	if !m.IsSensitive("api_key") {
		t.Error("expected case-insensitive match for api_key")
	}
	if m.IsSensitive("HOST") {
		t.Error("expected HOST to not be sensitive")
	}
}

func TestMasker_MaskFull(t *testing.T) {
	m := NewMasker(MaskFull, "SECRET")

	got := m.Mask("SECRET", "mysecretvalue")
	want := "*************"
	if got != want {
		t.Errorf("MaskFull: got %q, want %q", got, want)
	}
}

func TestMasker_MaskPartial(t *testing.T) {
	m := NewMasker(MaskPartial, "TOKEN")

	got := m.Mask("TOKEN", "abcdef")
	want := "a****f"
	if got != want {
		t.Errorf("MaskPartial: got %q, want %q", got, want)
	}
}

func TestMasker_MaskPartial_ShortValue(t *testing.T) {
	m := NewMasker(MaskPartial, "TOKEN")

	got := m.Mask("TOKEN", "ab")
	if got != "**" {
		t.Errorf("expected '**', got %q", got)
	}
}

func TestMasker_MaskHash(t *testing.T) {
	m := NewMasker(MaskHash, "PASSWORD")

	got := m.Mask("PASSWORD", "hunter2")
	if got != "[REDACTED]" {
		t.Errorf("MaskHash: got %q, want [REDACTED]", got)
	}
}

func TestMasker_MaskEmptyValue(t *testing.T) {
	m := NewMasker(MaskFull, "KEY")

	got := m.Mask("KEY", "")
	if got != "" {
		t.Errorf("expected empty string for empty value, got %q", got)
	}
}

func TestMasker_NonSensitivePassthrough(t *testing.T) {
	m := NewMasker(MaskFull, "SECRET")

	got := m.Mask("HOST", "localhost")
	if got != "localhost" {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestMasker_Apply(t *testing.T) {
	m := NewMasker(MaskHash, "API_KEY")

	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "super-secret"},
		{Key: "PORT", Value: "8080"},
	}

	result := m.Apply(entries)

	if result[0].Value != "localhost" {
		t.Errorf("HOST should be unchanged, got %q", result[0].Value)
	}
	if result[1].Value != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted, got %q", result[1].Value)
	}
	if result[2].Value != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result[2].Value)
	}

	// Ensure original is not mutated
	if entries[1].Value != "super-secret" {
		t.Error("Apply should not mutate original entries")
	}
}
