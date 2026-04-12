package env

import "fmt"

// SnapshotDiffResult holds the diff between two snapshots.
type SnapshotDiffResult struct {
	FromID string
	ToID   string
	Diff   DiffResult
}

// String returns a human-readable summary of the snapshot diff.
func (r *SnapshotDiffResult) String() string {
	if !HasDiff(r.Diff) {
		return fmt.Sprintf("snapshots %q and %q are identical", r.FromID, r.ToID)
	}
	return fmt.Sprintf("snapshots %q → %q: +%d -%d ~%d",
		r.FromID, r.ToID,
		len(r.Diff.Added),
		len(r.Diff.Removed),
		len(r.Diff.Changed),
	)
}

// CompareSnapshots diffs two snapshots by their IDs using the given manager.
func CompareSnapshots(sm *SnapshotManager, fromID, toID string) (*SnapshotDiffResult, error) {
	from, ok := sm.Get(fromID)
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", fromID)
	}
	to, ok := sm.Get(toID)
	if !ok {
		return nil, fmt.Errorf("snapshot %q not found", toID)
	}

	diffResult := Diff(from.Entries, to.Entries)
	return &SnapshotDiffResult{
		FromID: fromID,
		ToID:   toID,
		Diff:   diffResult,
	}, nil
}
