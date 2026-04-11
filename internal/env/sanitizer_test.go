package env

import (
	"testing"
)

func TestSanitizer_TrimKeysAndValues(t *testing.T) {
	s := NewSanitizer()
	input := []Entry{
		{Key: "  FOO  ", Value: "  bar  "},
		{Key: "BAZ", Value: "qux"},
	}
	out := s.Sanitize(input)
	if out[0].Key != "FOO" {
		t.Errorf("expected trimmed key FOO, got %q", out[0].Key)
	}
	if out[0].Value != "bar" {
		t.Errorf("expected trimmed value bar, got %q", out[0].Value)
	}
}

func TestSanitizer_StripQuotes(t *testing.T) {
	s := NewSanitizer()
	input := []Entry{
		{Key: "A", Value: `"hello world"`},
		{Key: "B", Value: "'single'"},
		{Key: "C", Value: "no-quotes"},
	}
	out := s.Sanitize(input)
	expected := []string{"hello world", "single", "no-quotes"}
	for i, e := range out {
		if e.Value != expected[i] {
			t.Errorf("entry %d: expected %q, got %q", i, expected[i], e.Value)
		}
	}
}

func TestSanitizer_WithoutStripQuotes(t *testing.T) {
	s := NewSanitizer(WithoutStripQuotes())
	input := []Entry{{Key: "A", Value: `"keep-me"`}}
	out := s.Sanitize(input)
	if out[0].Value != `"keep-me"` {
		t.Errorf("expected quotes preserved, got %q", out[0].Value)
	}
}

func TestSanitizer_NormalizeKey(t *testing.T) {
	s := NewSanitizer(WithNormalizeKey())
	input := []Entry{
		{Key: "my-key", Value: "v"},
		{Key: "hello world", Value: "v"},
		{Key: "already_UPPER", Value: "v"},
	}
	out := s.Sanitize(input)
	expected := []string{"MY_KEY", "HELLO_WORLD", "ALREADY_UPPER"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("entry %d: expected key %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestSanitizer_SkipsEmptyKeys(t *testing.T) {
	s := NewSanitizer()
	input := []Entry{
		{Key: "   ", Value: "orphan"},
		{Key: "VALID", Value: "ok"},
	}
	out := s.Sanitize(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "VALID" {
		t.Errorf("expected VALID, got %q", out[0].Key)
	}
}

func TestSanitizer_EmptyInput(t *testing.T) {
	s := NewSanitizer()
	out := s.Sanitize([]Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
