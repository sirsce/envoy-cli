package env

import "strings"

// StripResult holds the outcome of stripping a single entry.
type StripResult struct {
	Key     string
	OldVal  string
	Stripped bool
}

// Stripper removes or blanks out entries whose values match given patterns or prefixes.
type Stripper struct {
	prefixes []string
	contains []string
	blanks   bool
}

// NewStripper creates a Stripper. By default, matched entries are removed.
// Use WithBlankInstead to replace values with empty strings instead of removing.
func NewStripper(opts ...func(*Stripper)) *Stripper {
	s := &Stripper{}
	for _, o := range opts {
		o(s)
	}
	return s
}

// WithStripValuePrefix adds a value-prefix pattern; entries whose values start
// with the prefix will be stripped.
func WithStripValuePrefix(prefix string) func(*Stripper) {
	return func(s *Stripper) {
		s.prefixes = append(s.prefixes, prefix)
	}
}

// WithStripValueContains adds a substring pattern; entries whose values contain
// the substring will be stripped.
func WithStripValueContains(sub string) func(*Stripper) {
	return func(s *Stripper) {
		s.contains = append(s.contains, sub)
	}
}

// WithBlankInstead makes the stripper blank values instead of removing entries.
func WithBlankInstead() func(*Stripper) {
	return func(s *Stripper) { s.blanks = true }
}

// Strip processes entries and returns cleaned entries plus a result log.
func (s *Stripper) Strip(entries []Entry) ([]Entry, []StripResult) {
	out := make([]Entry, 0, len(entries))
	results := []StripResult{}

	for _, e := range entries {
		if s.matches(e.Value) {
			results = append(results, StripResult{Key: e.Key, OldVal: e.Value, Stripped: true})
			if s.blanks {
				out = append(out, Entry{Key: e.Key, Value: ""})
			}
			continue
		}
		out = append(out, e)
	}
	return out, results
}

// StrippedCount returns the number of results that were stripped.
func StrippedCount(results []StripResult) int {
	n := 0
	for _, r := range results {
		if r.Stripped {
			n++
		}
	}
	return n
}

func (s *Stripper) matches(val string) bool {
	for _, p := range s.prefixes {
		if strings.HasPrefix(val, p) {
			return true
		}
	}
	for _, c := range s.contains {
		if strings.Contains(val, c) {
			return true
		}
	}
	return false
}
