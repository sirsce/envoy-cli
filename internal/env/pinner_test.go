package env

import (
	"testing"
	"time"
)

var fixedPinTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestPinner() *Pinner {
	return NewPinner(WithPinnerClock(func() time.Time { return fixedPinTime }))
}

func TestPinner_PinAndIsPinned(t *testing.T) {
	p := newTestPinner()
	p.Pin("SECRET_KEY", "abc123", "locked for release")

	if !p.IsPinned("SECRET_KEY") {
		t.Fatal("expected SECRET_KEY to be pinned")
	}
	if p.IsPinned("OTHER_KEY") {
		t.Fatal("expected OTHER_KEY to not be pinned")
	}
}

func TestPinner_Get(t *testing.T) {
	p := newTestPinner()
	p.Pin("DB_PASS", "secret", "do not rotate")

	e, ok := p.Get("DB_PASS")
	if !ok {
		t.Fatal("expected to find pinned entry")
	}
	if e.Value != "secret" {
		t.Errorf("expected value %q, got %q", "secret", e.Value)
	}
	if e.Reason != "do not rotate" {
		t.Errorf("expected reason %q, got %q", "do not rotate", e.Reason)
	}
	if !e.PinnedAt.Equal(fixedPinTime) {
		t.Errorf("unexpected PinnedAt: %v", e.PinnedAt)
	}
}

func TestPinner_Unpin(t *testing.T) {
	p := newTestPinner()
	p.Pin("API_KEY", "val", "")

	if err := p.Unpin("API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.IsPinned("API_KEY") {
		t.Fatal("expected API_KEY to be unpinned")
	}
}

func TestPinner_Unpin_Missing(t *testing.T) {
	p := newTestPinner()
	if err := p.Unpin("NONEXISTENT"); err == nil {
		t.Fatal("expected error for unpinning missing key")
	}
}

func TestPinner_Apply(t *testing.T) {
	p := newTestPinner()
	p.Pin("SECRET", "pinned-value", "locked")

	entries := []Entry{
		{Key: "SECRET", Value: "incoming-value"},
		{Key: "OTHER", Value: "unchanged"},
	}

	result := p.Apply(entries)

	if result[0].Value != "pinned-value" {
		t.Errorf("expected pinned value, got %q", result[0].Value)
	}
	if result[1].Value != "unchanged" {
		t.Errorf("expected unchanged value, got %q", result[1].Value)
	}
	// original must not be mutated
	if entries[0].Value != "incoming-value" {
		t.Error("Apply must not mutate original entries")
	}
}

func TestPinner_List(t *testing.T) {
	p := newTestPinner()
	p.Pin("A", "1", "")
	p.Pin("B", "2", "")

	list := p.List()
	if len(list) != 2 {
		t.Errorf("expected 2 pinned entries, got %d", len(list))
	}
}
