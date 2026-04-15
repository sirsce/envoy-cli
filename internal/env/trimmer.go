package env

import (
	"strings"
)

// TrimMode controls which entries are affected by the Trimmer.
type TrimMode int

const (
	// TrimKeys trims whitespace from entry keys.
	TrimKeys TrimMode = 1 << iota
	// TrimValues trims whitespace from entry values.
	TrimValues
	// TrimBoth trims whitespace from both keys and values.
	TrimBoth = TrimKeys | TrimValues
)

// TrimResult holds the outcome of trimming a single entry.
type TrimResult struct {
	OriginalKey   string
	OriginalValue string
	NewKey        string
	NewValue      string
	Modified      bool
}

// Trimmer removes leading and trailing whitespace from entry keys and/or values.
type Trimmer struct {
	mode TrimMode
}

// NewTrimmer creates a new Trimmer with the given mode.
func NewTrimmer(mode TrimMode) *Trimmer {
	return &Trimmer{mode: mode}
}

// Trim applies whitespace trimming to the provided entries and returns
// the cleaned entries along with a result record for each entry.
func (t *Trimmer) Trim(entries []Entry) ([]Entry, []TrimResult) {
	out := make([]Entry, len(entries))
	results := make([]TrimResult, len(entries))

	for i, e := range entries {
		newKey := e.Key
		newVal := e.Value

		if t.mode&TrimKeys != 0 {
			newKey = strings.TrimSpace(e.Key)
		}
		if t.mode&TrimValues != 0 {
			newVal = strings.TrimSpace(e.Value)
		}

		modified := newKey != e.Key || newVal != e.Value
		results[i] = TrimResult{
			OriginalKey:   e.Key,
			OriginalValue: e.Value,
			NewKey:        newKey,
			NewValue:      newVal,
			Modified:      modified,
		}
		out[i] = Entry{Key: newKey, Value: newVal}
	}

	return out, results
}

// ModifiedCount returns the number of entries that were changed.
func ModifiedCount(results []TrimResult) int {
	count := 0
	for _, r := range results {
		if r.Modified {
			count++
		}
	}
	return count
}
