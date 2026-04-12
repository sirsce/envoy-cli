package env

import (
	"fmt"
	"sort"
)

// CompareResult holds the outcome of comparing two env entry sets.
type CompareResult struct {
	OnlyInLeft  []Entry
	OnlyInRight []Entry
	Changed     []EntryChange
	Identical   []Entry
}

// EntryChange represents a key whose value differs between two sets.
type EntryChange struct {
	Key      string
	OldValue string
	NewValue string
}

// Comparator compares two slices of Entry values.
type Comparator struct {
	ignoreCase bool
}

// ComparatorOption configures a Comparator.
type ComparatorOption func(*Comparator)

// WithCaseInsensitive makes value comparison case-insensitive.
func WithCaseInsensitive() ComparatorOption {
	return func(c *Comparator) {
		c.ignoreCase = true
	}
}

// NewComparator creates a new Comparator with optional options.
func NewComparator(opts ...ComparatorOption) *Comparator {
	c := &Comparator{}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Compare returns a CompareResult describing differences between left and right.
func (c *Comparator) Compare(left, right []Entry) CompareResult {
	leftMap := toMap(left)
	rightMap := toMap(right)

	result := CompareResult{}

	for _, e := range left {
		rVal, ok := rightMap[e.Key]
		if !ok {
			result.OnlyInLeft = append(result.OnlyInLeft, e)
			continue
		}
		if c.valuesEqual(e.Value, rVal) {
			result.Identical = append(result.Identical, e)
		} else {
			result.Changed = append(result.Changed, EntryChange{
				Key:      e.Key,
				OldValue: e.Value,
				NewValue: rVal,
			})
		}
	}

	for _, e := range right {
		if _, ok := leftMap[e.Key]; !ok {
			result.OnlyInRight = append(result.OnlyInRight, e)
		}
	}

	sortEntries(result.OnlyInLeft)
	sortEntries(result.OnlyInRight)
	sortEntries(result.Identical)
	sort.Slice(result.Changed, func(i, j int) bool {
		return result.Changed[i].Key < result.Changed[j].Key
	})

	return result
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	return fmt.Sprintf(
		"identical=%d changed=%d only_in_left=%d only_in_right=%d",
		len(r.Identical), len(r.Changed), len(r.OnlyInLeft), len(r.OnlyInRight),
	)
}

func (c *Comparator) valuesEqual(a, b string) bool {
	if c.ignoreCase {
		return stringsEqualFold(a, b)
	}
	return a == b
}

func stringsEqualFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 32
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 32
		}
		if ca != cb {
			return false
		}
	}
	return true
}

func toMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

func sortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
}
