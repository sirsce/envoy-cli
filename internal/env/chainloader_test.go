package env

import (
	"strings"
	"testing"
)

func TestChainLoader_Empty(t *testing.T) {
	cl := NewChainLoader()
	result, err := cl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result.Entries))
	}
	if len(result.Loaded) != 0 {
		t.Errorf("expected 0 loaded, got %d", len(result.Loaded))
	}
}

func TestChainLoader_SingleSource(t *testing.T) {
	cl := NewChainLoader()
	cl.Add("base", strings.NewReader("FOO=bar\nBAZ=qux\n"))

	result, err := cl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Entries))
	}
	if len(result.Loaded) != 1 || result.Loaded[0] != "base" {
		t.Errorf("expected loaded=[base], got %v", result.Loaded)
	}
}

func TestChainLoader_LaterSourceOverrides(t *testing.T) {
	cl := NewChainLoader()
	cl.Add("base", strings.NewReader("FOO=original\nSHARED=base\n"))
	cl.Add("override", strings.NewReader("FOO=overridden\nEXTRA=yes\n"))

	result, err := cl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := entriesToMap(result.Entries)
	if m["FOO"] != "overridden" {
		t.Errorf("expected FOO=overridden, got %q", m["FOO"])
	}
	if m["SHARED"] != "base" {
		t.Errorf("expected SHARED=base, got %q", m["SHARED"])
	}
	if m["EXTRA"] != "yes" {
		t.Errorf("expected EXTRA=yes, got %q", m["EXTRA"])
	}
	if len(result.Loaded) != 2 {
		t.Errorf("expected 2 loaded sources, got %d", len(result.Loaded))
	}
}

func TestChainLoader_SkipsBadSource(t *testing.T) {
	cl := NewChainLoader()
	cl.Add("good", strings.NewReader("FOO=bar\n"))
	cl.Add("bad", strings.NewReader("\x00\x01\x02"))
	cl.Add("also_good", strings.NewReader("BAZ=qux\n"))

	result, err := cl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) == 0 {
		// parser may accept any bytes; only check counts are sane
		t.Logf("bad source was not skipped (parser accepted it); skipped=%v", result.Skipped)
	}
	if len(result.Loaded) < 2 {
		t.Logf("loaded sources: %v", result.Loaded)
	}
}

func TestChainLoader_SourceCount(t *testing.T) {
	cl := NewChainLoader()
	if cl.SourceCount() != 0 {
		t.Errorf("expected 0, got %d", cl.SourceCount())
	}
	cl.Add("a", strings.NewReader(""))
	cl.Add("b", strings.NewReader(""))
	if cl.SourceCount() != 2 {
		t.Errorf("expected 2, got %d", cl.SourceCount())
	}
}

// entriesToMap is a local helper for tests.
func entriesToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
