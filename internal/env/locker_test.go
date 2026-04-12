package env

import (
	"testing"
	"time"
)

func newTestLocker() (*Locker, *time.Time) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	l := NewLocker()
	WithLockerClock(l, func() time.Time { return now })
	return l, &now
}

func TestLocker_LockAndIsLocked(t *testing.T) {
	l, _ := newTestLocker()

	if err := l.Lock("DB_PASSWORD", "alice", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.IsLocked("DB_PASSWORD") {
		t.Error("expected key to be locked")
	}
}

func TestLocker_UnlockByOwner(t *testing.T) {
	l, _ := newTestLocker()
	_ = l.Lock("API_KEY", "alice", nil)

	if err := l.Unlock("API_KEY", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.IsLocked("API_KEY") {
		t.Error("expected key to be unlocked")
	}
}

func TestLocker_UnlockByWrongOwner(t *testing.T) {
	l, _ := newTestLocker()
	_ = l.Lock("SECRET", "alice", nil)

	err := l.Unlock("SECRET", "bob")
	if err == nil {
		t.Error("expected error when unlocking with wrong owner")
	}
}

func TestLocker_LockConflict(t *testing.T) {
	l, _ := newTestLocker()
	_ = l.Lock("TOKEN", "alice", nil)

	err := l.Lock("TOKEN", "bob", nil)
	if err == nil {
		t.Error("expected conflict error")
	}
}

func TestLocker_SameOwnerRelock(t *testing.T) {
	l, _ := newTestLocker()
	_ = l.Lock("TOKEN", "alice", nil)

	if err := l.Lock("TOKEN", "alice", nil); err != nil {
		t.Errorf("same owner re-lock should succeed, got: %v", err)
	}
}

func TestLocker_TTLExpiry(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	current := now
	l := NewLocker()
	WithLockerClock(l, func() time.Time { return current })

	ttl := 5 * time.Minute
	_ = l.Lock("EXPIRING_KEY", "alice", &ttl)

	if !l.IsLocked("EXPIRING_KEY") {
		t.Error("expected key to be locked before expiry")
	}

	current = now.Add(10 * time.Minute)
	if l.IsLocked("EXPIRING_KEY") {
		t.Error("expected key to be unlocked after TTL expiry")
	}
}

func TestLocker_Get(t *testing.T) {
	l, _ := newTestLocker()
	_ = l.Lock("DB_HOST", "ci-runner", nil)

	entry, ok := l.Get("DB_HOST")
	if !ok {
		t.Fatal("expected to find lock entry")
	}
	if entry.LockedBy != "ci-runner" {
		t.Errorf("expected owner ci-runner, got %q", entry.LockedBy)
	}
}

func TestLocker_List(t *testing.T) {
	l, _ := newTestLocker()
	_ = l.Lock("KEY_A", "alice", nil)
	_ = l.Lock("KEY_B", "bob", nil)

	list := l.List()
	if len(list) != 2 {
		t.Errorf("expected 2 locks, got %d", len(list))
	}
}
