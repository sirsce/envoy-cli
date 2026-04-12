package env

import (
	"fmt"
	"time"
)

// Snapshot captures the state of env entries at a point in time.
type Snapshot struct {
	ID        string
	Label     string
	CreatedAt time.Time
	Entries   []Entry
}

// SnapshotManager manages named snapshots of env entry sets.
type SnapshotManager struct {
	clock     func() time.Time
	snapshots map[string]*Snapshot
}

// NewSnapshotManager creates a new SnapshotManager.
func NewSnapshotManager() *SnapshotManager {
	return &SnapshotManager{
		clock:     time.Now,
		snapshots: make(map[string]*Snapshot),
	}
}

// WithSnapshotClock overrides the clock used for snapshot timestamps.
func WithSnapshotClock(clock func() time.Time) func(*SnapshotManager) {
	return func(sm *SnapshotManager) {
		sm.clock = clock
	}
}

// Take creates a snapshot with the given label and entries.
func (sm *SnapshotManager) Take(label string, entries []Entry) *Snapshot {
	copied := make([]Entry, len(entries))
	copy(copied, entries)
	snap := &Snapshot{
		ID:        fmt.Sprintf("%s-%d", label, sm.clock().UnixNano()),
		Label:     label,
		CreatedAt: sm.clock(),
		Entries:   copied,
	}
	sm.snapshots[snap.ID] = snap
	return snap
}

// Get retrieves a snapshot by ID.
func (sm *SnapshotManager) Get(id string) (*Snapshot, bool) {
	snap, ok := sm.snapshots[id]
	return snap, ok
}

// List returns all stored snapshots in insertion-independent order.
func (sm *SnapshotManager) List() []*Snapshot {
	result := make([]*Snapshot, 0, len(sm.snapshots))
	for _, s := range sm.snapshots {
		result = append(result, s)
	}
	return result
}

// Delete removes a snapshot by ID. Returns false if not found.
func (sm *SnapshotManager) Delete(id string) bool {
	if _, ok := sm.snapshots[id]; !ok {
		return false
	}
	delete(sm.snapshots, id)
	return true
}

// Restore returns a copy of the entries from the snapshot with the given ID.
func (sm *SnapshotManager) Restore(id string) ([]Entry, error) {
	snap, ok := sm.snapshots[id]
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", id)
	}
	copied := make([]Entry, len(snap.Entries))
	copy(copied, snap.Entries)
	return copied, nil
}
