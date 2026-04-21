package env

import "strings"

// FlattenOptions configures how nested key structures are flattened.
type FlattenOptions struct {
	Separator string
	Prefix    string
	Uppercase bool
}

// Flattener collapses dot-notated or separator-delimited key hierarchies
// into a flat list of env entries, optionally applying a prefix.
type Flattener struct {
	opts FlattenOptions
}

// NewFlattener creates a Flattener with the given options.
// If Separator is empty, "." is used as the default.
func NewFlattener(opts FlattenOptions) *Flattener {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	return &Flattener{opts: opts}
}

// Flatten takes a slice of entries whose keys may contain the separator
// and returns a new slice with keys normalised to underscore-joined segments.
// An optional prefix is prepended to every resulting key.
func (f *Flattener) Flatten(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := f.normalise(e.Key)
		if f.opts.Prefix != "" {
			key = f.normalise(f.opts.Prefix) + "_" + key
		}
		out = append(out, Entry{Key: key, Value: e.Value})
	}
	return out
}

// FlattenedCount returns the number of keys that would be modified by Flatten.
func (f *Flattener) FlattenedCount(entries []Entry) int {
	count := 0
	for _, e := range entries {
		if strings.Contains(e.Key, f.opts.Separator) {
			count++
		}
	}
	return count
}

func (f *Flattener) normalise(key string) string {
	result := strings.ReplaceAll(key, f.opts.Separator, "_")
	if f.opts.Uppercase {
		result = strings.ToUpper(result)
	}
	return result
}
