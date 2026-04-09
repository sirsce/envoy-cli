package env

// DiffType represents the type of change between two env sets.
type DiffType string

const (
	DiffAdded   DiffType = "added"
	DiffRemoved DiffType = "removed"
	DiffChanged DiffType = "changed"
)

// DiffEntry represents a single difference between two env maps.
type DiffEntry struct {
	Key      string
	Type     DiffType
	OldValue string
	NewValue string
}

// Diff computes the differences between two env maps (base vs target).
// Returns a slice of DiffEntry describing what changed.
func Diff(base, target map[string]string) []DiffEntry {
	var diffs []DiffEntry

	// Check for changed or removed keys
	for k, baseVal := range base {
		targetVal, exists := target[k]
		if !exists {
			diffs = append(diffs, DiffEntry{
				Key:      k,
				Type:     DiffRemoved,
				OldValue: baseVal,
			})
		} else if baseVal != targetVal {
			diffs = append(diffs, DiffEntry{
				Key:      k,
				Type:     DiffChanged,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		}
	}

	// Check for added keys
	for k, targetVal := range target {
		if _, exists := base[k]; !exists {
			diffs = append(diffs, DiffEntry{
				Key:      k,
				Type:     DiffAdded,
				NewValue: targetVal,
			})
		}
	}

	return diffs
}

// HasDiff returns true if there are any differences between base and target.
func HasDiff(base, target map[string]string) bool {
	return len(Diff(base, target)) > 0
}
