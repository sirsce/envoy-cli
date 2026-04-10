package env

// MergeStrategy defines how conflicting keys are resolved during merge.
type MergeStrategy int

const (
	// MergeStrategyKeepBase keeps the base value on conflict.
	MergeStrategyKeepBase MergeStrategy = iota
	// MergeStrategyOverride replaces the base value with the override value.
	MergeStrategyOverride
	// MergeStrategyError returns an error on conflict.
	MergeStrategyError
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Entries  []Entry
	Added    []string
	Overridden []string
	Skipped  []string
}

// Merger merges two sets of env entries using a configurable strategy.
type Merger struct {
	strategy MergeStrategy
}

// NewMerger creates a new Merger with the given strategy.
func NewMerger(strategy MergeStrategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines base and override entries according to the merger's strategy.
// Returns a MergeResult and an error if the strategy is MergeStrategyError and
// a conflict is detected.
func (m *Merger) Merge(base, override []Entry) (MergeResult, error) {
	result := MergeResult{}
	index := make(map[string]int, len(base))
	merged := make([]Entry, len(base))
	copy(merged, base)

	for i, e := range merged {
		index[e.Key] = i
	}

	for _, e := range override {
		if idx, exists := index[e.Key]; exists {
			switch m.strategy {
			case MergeStrategyError:
				return MergeResult{}, &MergeConflictError{Key: e.Key}
			case MergeStrategyOverride:
				merged[idx] = e
				result.Overridden = append(result.Overridden, e.Key)
			case MergeStrategyKeepBase:
				result.Skipped = append(result.Skipped, e.Key)
			}
		} else {
			merged = append(merged, e)
			result.Added = append(result.Added, e.Key)
		}
	}

	result.Entries = merged
	return result, nil
}

// MergeConflictError is returned when a key conflict is detected under MergeStrategyError.
type MergeConflictError struct {
	Key string
}

func (e *MergeConflictError) Error() string {
	return "merge conflict on key: " + e.Key
}
