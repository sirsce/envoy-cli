package env

import (
	"bytes"
	"strings"
	"testing"
)

func makeInjectorEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestInjector_NewKeyAdded(t *testing.T) {
	inj := NewInjector(InjectorSkip)
	target := makeInjectorEntries("A", "1", "B", "2")
	source := makeInjectorEntries("C", "3")
	out, results, err := inj.Inject(target, source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if InjectedCount(results) != 1 {
		t.Errorf("expected 1 injected, got %d", InjectedCount(results))
	}
}

func TestInjector_SkipExisting(t *testing.T) {
	inj := NewInjector(InjectorSkip)
	target := makeInjectorEntries("A", "1")
	source := makeInjectorEntries("A", "99")
	out, results, err := inj.Inject(target, source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "1" {
		t.Errorf("expected original value '1', got %q", out[0].Value)
	}
	if results[0].Skipped != true {
		t.Error("expected result to be skipped")
	}
}

func TestInjector_OverwriteExisting(t *testing.T) {
	inj := NewInjector(InjectorOverwrite)
	target := makeInjectorEntries("A", "old")
	source := makeInjectorEntries("A", "new")
	out, results, err := inj.Inject(target, source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "new" {
		t.Errorf("expected overwritten value 'new', got %q", out[0].Value)
	}
	if InjectedCount(results) != 1 {
		t.Errorf("expected 1 injected, got %d", InjectedCount(results))
	}
}

func TestInjector_ErrorOnConflict(t *testing.T) {
	inj := NewInjector(InjectorError)
	target := makeInjectorEntries("X", "1")
	source := makeInjectorEntries("X", "2")
	_, _, err := inj.Inject(target, source)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
	if !strings.Contains(err.Error(), "X") {
		t.Errorf("error should mention key X: %v", err)
	}
}

func TestInjector_DoesNotMutateTarget(t *testing.T) {
	inj := NewInjector(InjectorOverwrite)
	target := makeInjectorEntries("A", "1")
	original := make([]Entry, len(target))
	copy(original, target)
	source := makeInjectorEntries("A", "modified")
	inj.Inject(target, source) //nolint:errcheck
	if target[0].Value != original[0].Value {
		t.Error("Inject mutated the original target slice")
	}
}

func TestInjectorReport_Text(t *testing.T) {
	results := []InjectorResult{
		{Key: "FOO", Injected: true, Reason: "new key"},
		{Key: "BAR", Skipped: true, Reason: "key already exists"},
	}
	rep := NewInjectorReport(results)
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO") || !strings.Contains(out, "BAR") {
		t.Errorf("unexpected report output: %s", out)
	}
}

func TestInjectorReport_JSON(t *testing.T) {
	results := []InjectorResult{
		{Key: "K", Injected: true, Reason: "new key"},
	}
	rep := NewInjectorReport(results)
	var buf bytes.Buffer
	if err := rep.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(buf.String(), `"key"`) {
		t.Errorf("expected JSON with 'key' field, got: %s", buf.String())
	}
}
