package env

import (
	"testing"
)

func makeNormalizerEntries() []Entry {
	return []Entry{
		{Key: "db_host", Value: "localhost"},
		{Key: "API-KEY", Value: "secret"},
		{Key: "AppName", Value: "envoy"},
	}
}

func TestNormalizer_UpperMode(t *testing.T) {
	n := NewNormalizer(WithNormalizeMode(NormalizeModeUpper))
	out := n.Normalize(makeNormalizerEntries())

	expected := []string{"DB_HOST", "API-KEY", "APPNAME"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestNormalizer_SnakeMode(t *testing.T) {
	n := NewNormalizer(WithNormalizeMode(NormalizeModeSnake))
	entries := []Entry{
		{Key: "DB-HOST", Value: "localhost"},
		{Key: "App Name", Value: "envoy"},
	}
	out := n.Normalize(entries)

	expected := []string{"db_host", "app_name"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestNormalizer_CamelMode(t *testing.T) {
	n := NewNormalizer(WithNormalizeMode(NormalizeModeCamel))
	entries := []Entry{
		{Key: "db_host", Value: "localhost"},
		{Key: "api_secret_key", Value: "abc"},
	}
	out := n.Normalize(entries)

	expected := []string{"dbHost", "apiSecretKey"}
	for i, e := range out {
		if e.Key != expected[i] {
			t.Errorf("entry %d: got key %q, want %q", i, e.Key, expected[i])
		}
	}
}

func TestNormalizer_ModifiedCount(t *testing.T) {
	n := NewNormalizer(WithNormalizeMode(NormalizeModeUpper))
	entries := []Entry{
		{Key: "ALREADY_UPPER", Value: "v1"},
		{Key: "needs_upper", Value: "v2"},
		{Key: "ALSO_UPPER", Value: "v3"},
	}
	n.Normalize(entries)
	if n.ModifiedCount() != 1 {
		t.Errorf("expected 1 modified, got %d", n.ModifiedCount())
	}
}

func TestNormalizer_DoesNotMutateOriginal(t *testing.T) {
	n := NewNormalizer()
	original := makeNormalizerEntries()
	copy := make([]Entry, len(original))
	for i, e := range original {
		copy[i] = e
	}
	n.Normalize(original)
	for i, e := range original {
		if e.Key != copy[i].Key {
			t.Errorf("entry %d mutated: got %q, want %q", i, e.Key, copy[i].Key)
		}
	}
}

func TestNormalizer_DefaultIsUpper(t *testing.T) {
	n := NewNormalizer()
	out := n.Normalize([]Entry{{Key: "lower", Value: "x"}})
	if out[0].Key != "LOWER" {
		t.Errorf("expected LOWER, got %q", out[0].Key)
	}
}
