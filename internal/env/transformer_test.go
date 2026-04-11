package env

import (
	"testing"
)

func TestTransformer_TrimSpace(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "  hello  "},
		{Key: "B", Value: "\tworld\n"},
	}
	tr := NewTransformer(TrimSpaceTransform())
	result := tr.Apply(entries)

	if result[0].Value != "hello" {
		t.Errorf("expected 'hello', got %q", result[0].Value)
	}
	if result[1].Value != "world" {
		t.Errorf("expected 'world', got %q", result[1].Value)
	}
}

func TestTransformer_ToUpper(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "lowercase"}}
	tr := NewTransformer(ToUpperTransform())
	result := tr.Apply(entries)
	if result[0].Value != "LOWERCASE" {
		t.Errorf("expected 'LOWERCASE', got %q", result[0].Value)
	}
}

func TestTransformer_ToLower(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "UPPER"}}
	tr := NewTransformer(ToLowerTransform())
	result := tr.Apply(entries)
	if result[0].Value != "upper" {
		t.Errorf("expected 'upper', got %q", result[0].Value)
	}
}

func TestTransformer_Replace(t *testing.T) {
	entries := []Entry{{Key: "URL", Value: "http://localhost:8080"}}
	tr := NewTransformer(ReplaceTransform("localhost", "example.com"))
	result := tr.Apply(entries)
	if result[0].Value != "http://example.com:8080" {
		t.Errorf("unexpected value: %q", result[0].Value)
	}
}

func TestTransformer_PrefixValue(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "world"},
		{Key: "B", Value: ""},
	}
	tr := NewTransformer(PrefixValueTransform("hello_"))
	result := tr.Apply(entries)
	if result[0].Value != "hello_world" {
		t.Errorf("expected 'hello_world', got %q", result[0].Value)
	}
	if result[1].Value != "" {
		t.Errorf("expected empty string for empty value, got %q", result[1].Value)
	}
}

func TestTransformer_Chain(t *testing.T) {
	entries := []Entry{{Key: "K", Value: "  hello world  "}}
	tr := NewTransformer(TrimSpaceTransform(), ToUpperTransform())
	result := tr.Apply(entries)
	if result[0].Value != "HELLO WORLD" {
		t.Errorf("expected 'HELLO WORLD', got %q", result[0].Value)
	}
}

func TestTransformer_DoesNotMutateOriginal(t *testing.T) {
	original := []Entry{{Key: "K", Value: "original"}}
	tr := NewTransformer(ToUpperTransform())
	_ = tr.Apply(original)
	if original[0].Value != "original" {
		t.Errorf("original entry was mutated")
	}
}

func TestTransformer_Empty(t *testing.T) {
	tr := NewTransformer(TrimSpaceTransform())
	result := tr.Apply([]Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result")
	}
}
