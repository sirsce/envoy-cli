package env

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// LockEntry represents a lock on an environment key.
type LockEntry struct {
	Key       string
	LockedBy  string
	LockedAt  time.Time
	ExpiresAt *time.Time
}

// Locker manages exclusive locks on env keys to prevent concurrent modification.
type Locker struct {
	mu    sync.RWMutex
	locks map[string]LockEntry
	clock func() time.Time
}

// NewLocker creates a new Locker instance.
func NewLocker() *Locker {
	return &Locker{
		locks: make(map[string]LockEntry),
		clock: time.Now,
	}
}

// WithLockerClock sets a custom clock for testing.
func WithLockerClock(l *Locker, clock func() time.Time) *Locker {
	l.clock = clock
	return l
}

// Lock acquires a lock on the given key for the specified owner.
// Returns an error if the key is already locked by another owner.
func (l *Locker) Lock(key, owner string, ttl *time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	if existing, ok := l.locks[key]; ok {
		if existing.ExpiresAt == nil || existing.ExpiresAt.After(now) {
			if existing.LockedBy != owner {
				return fmt.Errorf("key %q is locked by %q", key, existing.LockedBy)
			}
			return nil
		}
	}

	entry := LockEntry{
		Key:      key,
		LockedBy: owner,
		LockedAt: now,
	}
	if ttl != nil {
		exp := now.Add(*ttl)
		entry.ExpiresAt = &exp
	}
	l.locks[key] = entry
	return nil
}

// Unlock releases the lock on the given key if owned by the given owner.
func (l *Locker) Unlock(key, owner string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	existing, ok := l.locks[key]
	if !ok {
		return errors.New("key is not locked")
	}
	if existing.LockedBy != owner {
		return fmt.Errorf("key %q is owned by %q, not %q", key, existing.LockedBy, owner)
	}
	delete(l.locks, key)
	return nil
}

// IsLocked returns true if the key is currently locked (and not expired).
func (l *Locker) IsLocked(key string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entry, ok := l.locks[key]
	if !ok {
		return false
	}
	if entry.ExpiresAt != nil && !entry.ExpiresAt.After(l.clock()) {
		return false
	}
	return true
}

// Get returns the LockEntry for the given key, if it exists and is not expired.
func (l *Locker) Get(key string) (LockEntry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entry, ok := l.locks[key]
	if !ok {
		return LockEntry{}, false
	}
	if entry.ExpiresAt != nil && !entry.ExpiresAt.After(l.clock()) {
		return LockEntry{}, false
	}
	return entry, true
}

// List returns all currently active (non-expired) lock entries.
func (l *Locker) List() []LockEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	now := l.clock()
	var result []LockEntry
	for _, entry := range l.locks {
		if entry.ExpiresAt == nil || entry.ExpiresAt.After(now) {
			result = append(result, entry)
		}
	}
	return result
}
