package env

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var fixedRotationTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func newTestRotator() *Rotator {
	return NewRotator(WithClock(func() time.Time { return fixedRotationTime }))
}

func TestRotator_NoRegisteredKeys(t *testing.T) {
	r := newTestRotator()
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	out, records, err := r.Rotate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected 0 records, got %d", len(records))
	}
	if out[0].Value != "bar" {
		t.Errorf("expected value unchanged, got %q", out[0].Value)
	}
}

func TestRotator_RotatesRegisteredKey(t *testing.T) {
	r := newTestRotator()
	r.Register("SECRET", func(old string) (string, error) {
		return "new-" + old, nil
	})

	entries := []Entry{{Key: "SECRET", Value: "abc"}}
	out, records, err := r.Rotate(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].OldValue != "abc" || records[0].NewValue != "new-abc" {
		t.Errorf("record mismatch: %+v", records[0])
	}
	if records[0].RotatedAt != fixedRotationTime {
		t.Errorf("unexpected rotation time: %v", records[0].RotatedAt)
	}
	if out[0].Value != "new-abc" {
		t.Errorf("expected new-abc, got %q", out[0].Value)
	}
}

func TestRotator_DoesNotMutateOriginal(t *testing.T) {
	r := newTestRotator()
	r.Register("TOKEN", func(_ string) (string, error) { return "rotated", nil })

	orig := []Entry{{Key: "TOKEN", Value: "original"}}
	r.Rotate(orig) //nolint
	if orig[0].Value != "original" {
		t.Error("original entries were mutated")
	}
}

func TestRotator_GeneratorError(t *testing.T) {
	r := newTestRotator()
	r.Register("KEY", func(_ string) (string, error) {
		return "", errors.New("generator failed")
	})

	_, _, err := r.Rotate([]Entry{{Key: "KEY", Value: "v"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRotator_Keys(t *testing.T) {
	r := newTestRotator()
	for _, k := range []string{"A", "B", "C"} {
		k := k
		r.Register(k, func(_ string) (string, error) { return fmt.Sprintf("val-%s", k), nil })
	}
	if len(r.Keys()) != 3 {
		t.Errorf("expected 3 keys, got %d", len(r.Keys()))
	}
}
