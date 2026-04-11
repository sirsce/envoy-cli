package env

import (
	"testing"
	"time"
)

func fixedTime() time.Time {
	return time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
}

func newTestAuditor() *Auditor {
	a := NewAuditor()
	a.now = fixedTime
	return a
}

func TestAuditor_NoChanges(t *testing.T) {
	a := newTestAuditor()
	base := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	events := a.Audit(base, base)
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestAuditor_AddedKey(t *testing.T) {
	a := newTestAuditor()
	base := []Entry{{Key: "FOO", Value: "bar"}}
	updated := []Entry{{Key: "FOO", Value: "bar"}, {Key: "NEW", Value: "val"}}
	events := a.Audit(base, updated)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Action != AuditActionAdd || events[0].Key != "NEW" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestAuditor_RemovedKey(t *testing.T) {
	a := newTestAuditor()
	base := []Entry{{Key: "FOO", Value: "bar"}, {Key: "OLD", Value: "val"}}
	updated := []Entry{{Key: "FOO", Value: "bar"}}
	events := a.Audit(base, updated)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Action != AuditActionRemove || events[0].Key != "OLD" || events[0].OldValue != "val" {
		t.Errorf("unexpected event: %+v", events[0])
	}
}

func TestAuditor_ChangedKey(t *testing.T) {
	a := newTestAuditor()
	base := []Entry{{Key: "FOO", Value: "old"}}
	updated := []Entry{{Key: "FOO", Value: "new"}}
	events := a.Audit(base, updated)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ev := events[0]
	if ev.Action != AuditActionChange || ev.OldValue != "old" || ev.NewValue != "new" {
		t.Errorf("unexpected event: %+v", ev)
	}
}

func TestAuditor_EventString(t *testing.T) {
	ts := fixedTime()
	cases := []struct {
		ev   AuditEvent
		want string
	}{
		{AuditEvent{Timestamp: ts, Action: AuditActionAdd, Key: "K", NewValue: "v"}, "[2024-01-15T12:00:00Z] ADD K=\"v\""},
		{AuditEvent{Timestamp: ts, Action: AuditActionRemove, Key: "K", OldValue: "v"}, "[2024-01-15T12:00:00Z] REMOVE K (was \"v\")"},
		{AuditEvent{Timestamp: ts, Action: AuditActionChange, Key: "K", OldValue: "a", NewValue: "b"}, "[2024-01-15T12:00:00Z] CHANGE K: \"a\" -> \"b\""},
	}
	for _, c := range cases {
		if got := c.ev.String(); got != c.want {
			t.Errorf("String() = %q, want %q", got, c.want)
		}
	}
}
