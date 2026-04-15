package env

import (
	"testing"
)

func makeTrimmerEntries() []Entry {
	return []Entry{
		{Key: "  KEY_A  ", Value: "  value_a  "},
		{Key: "KEY_B", Value: "value_b"},
		{Key: " KEY_C", Value: "value_c "},
		{Key: "KEY_D", Value: "  "},
	}
}

func TestTrimmer_TrimBoth(t *testing.T) {
	tr := NewTrimmer(TrimBoth)
	out, results := tr.Trim(makeTrimmerEntries())

	if out[0].Key != "KEY_A" {
		t.Errorf("expected KEY_A, got %q", out[0].Key)
	}
	if out[0].Value != "value_a" {
		t.Errorf("expected value_a, got %q", out[0].Value)
	}
	if !results[0].Modified {
		t.Error("expected entry 0 to be marked modified")
	}
	if results[1].Modified {
		t.Error("expected entry 1 to be unmodified")
	}
}

func TestTrimmer_TrimKeysOnly(t *testing.T) {
	tr := NewTrimmer(TrimKeys)
	out, results := tr.Trim(makeTrimmerEntries())

	if out[0].Key != "KEY_A" {
		t.Errorf("expected KEY_A, got %q", out[0].Key)
	}
	// value should be unchanged
	if out[0].Value != "  value_a  " {
		t.Errorf("expected original value, got %q", out[0].Value)
	}
	if !results[0].Modified {
		t.Error("expected entry 0 to be modified (key changed)")
	}
}

func TestTrimmer_TrimValuesOnly(t *testing.T) {
	tr := NewTrimmer(TrimValues)
	out, results := tr.Trim(makeTrimmerEntries())

	// key should be unchanged
	if out[0].Key != "  KEY_A  " {
		t.Errorf("expected original key, got %q", out[0].Key)
	}
	if out[0].Value != "value_a" {
		t.Errorf("expected value_a, got %q", out[0].Value)
	}
	if !results[0].Modified {
		t.Error("expected entry 0 to be modified (value changed)")
	}
}

func TestTrimmer_DoesNotMutateOriginal(t *testing.T) {
	original := makeTrimmerEntries()
	tr := NewTrimmer(TrimBoth)
	tr.Trim(original)

	if original[0].Key != "  KEY_A  " {
		t.Error("original entries should not be mutated")
	}
}

func TestTrimmer_ModifiedCount(t *testing.T) {
	tr := NewTrimmer(TrimBoth)
	_, results := tr.Trim(makeTrimmerEntries())

	count := ModifiedCount(results)
	// entries 0, 2, 3 are modified (KEY_B/value_b are already clean)
	if count != 3 {
		t.Errorf("expected 3 modified, got %d", count)
	}
}

func TestTrimmer_EmptyEntries(t *testing.T) {
	tr := NewTrimmer(TrimBoth)
	out, results := tr.Trim([]Entry{})

	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
