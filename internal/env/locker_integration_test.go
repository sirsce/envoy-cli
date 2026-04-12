package env

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestLocker_Integration_LockAndReport(t *testing.T) {
	now := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	l := NewLocker()
	WithLockerClock(l, func() time.Time { return now })

	keys := []string{"DB_PASS", "API_SECRET", "JWT_KEY"}
	for _, k := range keys {
		if err := l.Lock(k, "deploy-bot", nil); err != nil {
			t.Fatalf("failed to lock %s: %v", k, err)
		}
	}

	list := l.List()
	report := NewLockReport(list)

	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}

	out := buf.String()
	for _, k := range keys {
		if !strings.Contains(out, k) {
			t.Errorf("expected %s in report output", k)
		}
	}
	if !strings.Contains(out, "deploy-bot") {
		t.Error("expected deploy-bot in report output")
	}
}

func TestLocker_Integration_UnlockReducesList(t *testing.T) {
	l := NewLocker()
	_ = l.Lock("KEY_A", "alice", nil)
	_ = l.Lock("KEY_B", "alice", nil)
	_ = l.Lock("KEY_C", "alice", nil)

	_ = l.Unlock("KEY_B", "alice")

	list := l.List()
	if len(list) != 2 {
		t.Errorf("expected 2 active locks after unlock, got %d", len(list))
	}
	for _, e := range list {
		if e.Key == "KEY_B" {
			t.Error("KEY_B should not appear in list after unlock")
		}
	}
}

func TestLocker_Integration_ExpiredLockAllowsRelock(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	current := now
	l := NewLocker()
	WithLockerClock(l, func() time.Time { return current })

	ttl := 2 * time.Minute
	_ = l.Lock("TEMP_KEY", "alice", &ttl)

	current = now.Add(5 * time.Minute)

	if err := l.Lock("TEMP_KEY", "bob", nil); err != nil {
		t.Errorf("expected expired lock to allow re-lock by new owner, got: %v", err)
	}

	entry, ok := l.Get("TEMP_KEY")
	if !ok {
		t.Fatal("expected to find re-locked entry")
	}
	if entry.LockedBy != "bob" {
		t.Errorf("expected bob to own lock, got %q", entry.LockedBy)
	}
}
