package env

import (
	"sort"
	"strings"
)

// SortOrder defines how entries should be sorted.
type SortOrder int

const (
	// SortAlpha sorts entries alphabetically by key.
	SortAlpha SortOrder = iota
	// SortAlphaDesc sorts entries reverse-alphabetically by key.
	SortAlphaDesc
	// SortByGroup sorts entries by group prefix (e.g. DB_, APP_) then alphabetically.
	SortByGroup
)

// Sorter sorts env entries by a given strategy.
type Sorter struct {
	order SortOrder
}

// NewSorter creates a new Sorter with the given order.
func NewSorter(order SortOrder) *Sorter {
	return &Sorter{order: order}
}

// Sort returns a new slice of entries sorted according to the Sorter's order.
func (s *Sorter) Sort(entries []Entry) []Entry {
	result := make([]Entry, len(entries))
	copy(result, entries)

	switch s.order {
	case SortAlpha:
		sort.Slice(result, func(i, j int) bool {
			return result[i].Key < result[j].Key
		})
	case SortAlphaDesc:
		sort.Slice(result, func(i, j int) bool {
			return result[i].Key > result[j].Key
		})
	case SortByGroup:
		sort.Slice(result, func(i, j int) bool {
			gi := groupPrefix(result[i].Key)
			gj := groupPrefix(result[j].Key)
			if gi != gj {
				return gi < gj
			}
			return result[i].Key < result[j].Key
		})
	}

	return result
}

// groupPrefix returns the prefix of a key up to and including the first underscore.
// If no underscore is found, the full key is returned.
func groupPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx != -1 {
		return key[:idx+1]
	}
	return key
}
