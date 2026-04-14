package env

import "fmt"

// DeduplicateStrategy controls how duplicate keys are resolved.
type DeduplicateStrategy int

const (
	// DeduplicateKeepFirst retains the first occurrence of a duplicate key.
	DeduplicateKeepFirst DeduplicateStrategy = iota
	// DeduplicateKeepLast retains the last occurrence of a duplicate key.
	DeduplicateKeepLast
	// DeduplicateError returns an error when a duplicate key is found.
	DeduplicateError
)

// DuplicateRecord describes a duplicate key that was found.
type DuplicateRecord struct {
	Key   string
	Count int
}

// DeduplicateResult holds the deduplicated entries and any duplicate records.
type DeduplicateResult struct {
	Entries    []Entry
	Duplicates []DuplicateRecord
}

// Deduplicator removes or reports duplicate keys from a slice of entries.
type Deduplicator struct {
	strategy DeduplicateStrategy
}

// NewDeduplicator creates a Deduplicator with the given strategy.
func NewDeduplicator(strategy DeduplicateStrategy) *Deduplicator {
	return &Deduplicator{strategy: strategy}
}

// Deduplicate processes entries and resolves duplicates according to the strategy.
func (d *Deduplicator) Deduplicate(entries []Entry) (DeduplicateResult, error) {
	seen := make(map[string]int) // key -> count
	order := make([]string, 0, len(entries))
	last := make(map[string]Entry)
	first := make(map[string]Entry)

	for _, e := range entries {
		if _, exists := seen[e.Key]; !exists {
			order = append(order, e.Key)
			first[e.Key] = e
		}
		seen[e.Key]++
		last[e.Key] = e
	}

	var duplicates []DuplicateRecord
	for _, key := range order {
		if seen[key] > 1 {
			duplicates = append(duplicates, DuplicateRecord{Key: key, Count: seen[key]})
		}
	}

	if d.strategy == DeduplicateError && len(duplicates) > 0 {
		return DeduplicateResult{}, fmt.Errorf("duplicate key found: %q (count: %d)", duplicates[0].Key, duplicates[0].Count)
	}

	result := make([]Entry, 0, len(order))
	for _, key := range order {
		if d.strategy == DeduplicateKeepLast {
			result = append(result, last[key])
		} else {
			result = append(result, first[key])
		}
	}

	return DeduplicateResult{Entries: result, Duplicates: duplicates}, nil
}

// HasDuplicates returns true if any duplicate keys exist in the provided entries.
func HasDuplicates(entries []Entry) bool {
	seen := make(map[string]struct{}, len(entries))
	for _, e := range entries {
		if _, ok := seen[e.Key]; ok {
			return true
		}
		seen[e.Key] = struct{}{}
	}
	return false
}
