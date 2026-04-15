package env

import (
	"testing"
)

func TestPlaceholderResolver_NoPlaceholders(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, true)
	entries := []Entry{{Key: "HOST", Value: "localhost"}, {Key: "PORT", Value: "8080"}}
	out, err := r.Resolve(entries, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" || out[1].Value != "8080" {
		t.Errorf("values mutated unexpectedly")
	}
}

func TestPlaceholderResolver_DoubleBrace(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, true)
	entries := []Entry{{Key: "DSN", Value: "postgres://{{ USER }}:{{ PASS }}@localhost/db"}}
	lookup := map[string]string{"USER": "admin", "PASS": "secret"}
	out, err := r.Resolve(entries, lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://admin:secret@localhost/db"
	if out[0].Value != want {
		t.Errorf("got %q, want %q", out[0].Value, want)
	}
}

func TestPlaceholderResolver_AngleBracket(t *testing.T) {
	r := NewPlaceholderResolver(StyleAngleBracket, true)
	entries := []Entry{{Key: "URL", Value: "http://<HOST>:<PORT>"}}
	lookup := map[string]string{"HOST": "example.com", "PORT": "443"}
	out, err := r.Resolve(entries, lookup)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "http://example.com:443" {
		t.Errorf("unexpected value: %s", out[0].Value)
	}
}

func TestPlaceholderResolver_PercentStyle(t *testing.T) {
	r := NewPlaceholderResolver(StylePercent, false)
	entries := []Entry{{Key: "GREETING", Value: "Hello %NAME%!"}}
	out, err := r.Resolve(entries, map[string]string{"NAME": "World"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "Hello World!" {
		t.Errorf("unexpected value: %s", out[0].Value)
	}
}

func TestPlaceholderResolver_StrictMissingError(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, true)
	entries := []Entry{{Key: "DSN", Value: "postgres://{{ MISSING }}@localhost"}}
	_, err := r.Resolve(entries, map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing placeholder, got nil")
	}
}

func TestPlaceholderResolver_NonStrictLeavesUnresolved(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, false)
	entries := []Entry{{Key: "X", Value: "{{ MISSING }}"}}
	out, err := r.Resolve(entries, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "{{ MISSING }}" {
		t.Errorf("expected original value, got %q", out[0].Value)
	}
}

func TestPlaceholderResolver_DetectPlaceholders(t *testing.T) {
	r := NewPlaceholderResolver(StyleDoubleBrace, false)
	keys := r.DetectPlaceholders("{{ A }} and {{ B }}")
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
