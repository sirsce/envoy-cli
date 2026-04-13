package env

import (
	"testing"
	"time"
)

var fixedArchiveTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestArchiver() *Archiver {
	return NewArchiver(WithArchiverClock(func() time.Time { return fixedArchiveTime }))
}

func TestArchiver_ArchiveAndGet(t *testing.T) {
	a := newTestArchiver()
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	v := a.Archive("initial", entries)
	if v != 1 {
		t.Fatalf("expected version 1, got %d", v)
	}
	ae, err := a.Get(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ae.Label != "initial" {
		t.Errorf("expected label 'initial', got %q", ae.Label)
	}
	if len(ae.Entries) != 1 || ae.Entries[0].Key != "FOO" {
		t.Errorf("unexpected entries: %+v", ae.Entries)
	}
	if !ae.ArchivedAt.Equal(fixedArchiveTime) {
		t.Errorf("unexpected timestamp: %v", ae.ArchivedAt)
	}
}

func TestArchiver_DoesNotMutateOriginal(t *testing.T) {
	a := newTestArchiver()
	entries := []Entry{{Key: "X", Value: "1"}}
	a.Archive("v1", entries)
	entries[0].Value = "mutated"
	ae, _ := a.Get(1)
	if ae.Entries[0].Value == "mutated" {
		t.Error("archiver should not share slice with caller")
	}
}

func TestArchiver_List(t *testing.T) {
	a := newTestArchiver()
	a.Archive("first", []Entry{{Key: "A", Value: "1"}})
	a.Archive("second", []Entry{{Key: "B", Value: "2"}})
	list := a.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 archives, got %d", len(list))
	}
	if list[0].Version != 1 || list[1].Version != 2 {
		t.Errorf("unexpected order: %+v", list)
	}
}

func TestArchiver_Delete(t *testing.T) {
	a := newTestArchiver()
	a.Archive("v1", []Entry{{Key: "A", Value: "1"}})
	a.Archive("v2", []Entry{{Key: "B", Value: "2"}})
	if err := a.Delete(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := a.Get(1); err == nil {
		t.Error("expected error after delete")
	}
	if len(a.List()) != 1 {
		t.Errorf("expected 1 remaining archive")
	}
}

func TestArchiver_Delete_NotFound(t *testing.T) {
	a := newTestArchiver()
	if err := a.Delete(99); err == nil {
		t.Error("expected error for missing version")
	}
}

func TestArchiver_Latest(t *testing.T) {
	a := newTestArchiver()
	_, err := a.Latest()
	if err == nil {
		t.Error("expected error when no archives")
	}
	a.Archive("v1", []Entry{{Key: "A", Value: "1"}})
	a.Archive("v2", []Entry{{Key: "B", Value: "2"}})
	latest, err := a.Latest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest.Version != 2 {
		t.Errorf("expected latest version 2, got %d", latest.Version)
	}
}
