package env

import "strings"

// FilterOptions controls how entries are filtered.
type FilterOptions struct {
	Prefix    string
	Suffix    string
	KeySubstr string
	ExcludeKeys []string
}

// Filter applies filtering logic to a slice of Entry values.
type Filter struct {
	opts FilterOptions
	excludeSet map[string]struct{}
}

// NewFilter creates a new Filter with the given options.
func NewFilter(opts FilterOptions) *Filter {
	exclude := make(map[string]struct{}, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		exclude[k] = struct{}{}
	}
	return &Filter{opts: opts, excludeSet: exclude}
}

// Apply returns only the entries that match the filter criteria.
func (f *Filter) Apply(entries []Entry) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !f.matches(e) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// matches returns true if the entry satisfies all filter conditions.
func (f *Filter) matches(e Entry) bool {
	if _, excluded := f.excludeSet[e.Key]; excluded {
		return false
	}
	if f.opts.Prefix != "" && !strings.HasPrefix(e.Key, f.opts.Prefix) {
		return false
	}
	if f.opts.Suffix != "" && !strings.HasSuffix(e.Key, f.opts.Suffix) {
		return false
	}
	if f.opts.KeySubstr != "" && !strings.Contains(e.Key, f.opts.KeySubstr) {
		return false
	}
	return true
}

// FilterByPrefix is a convenience function that filters entries by key prefix.
func FilterByPrefix(entries []Entry, prefix string) []Entry {
	return NewFilter(FilterOptions{Prefix: prefix}).Apply(entries)
}

// StripPrefix returns entries with the given prefix removed from each key.
func StripPrefix(entries []Entry, prefix string) []Entry {
	result := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if strings.HasPrefix(e.Key, prefix) {
			e.Key = strings.TrimPrefix(e.Key, prefix)
		}
		result = append(result, e)
	}
	return result
}

// FilterAndStrip is a convenience function that filters entries by key prefix
// and then strips that prefix from each matching entry's key.
func FilterAndStrip(entries []Entry, prefix string) []Entry {
	return StripPrefix(FilterByPrefix(entries, prefix), prefix)
}
