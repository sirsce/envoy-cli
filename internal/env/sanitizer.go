package env

import (
	"strings"
	"unicode"
)

// SanitizeOption configures sanitization behavior.
type SanitizeOption func(*Sanitizer)

// Sanitizer cleans env entry keys and values.
type Sanitizer struct {
	trimKeys    bool
	trimValues  bool
	stripQuotes bool
	normalizeKey bool
}

// NewSanitizer creates a Sanitizer with the given options.
func NewSanitizer(opts ...SanitizeOption) *Sanitizer {
	s := &Sanitizer{
		trimKeys:    true,
		trimValues:  true,
		stripQuotes: true,
		normalizeKey: false,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// WithNormalizeKey enables uppercasing and replacing spaces/hyphens with underscores in keys.
func WithNormalizeKey() SanitizeOption {
	return func(s *Sanitizer) { s.normalizeKey = true }
}

// WithoutStripQuotes disables stripping surrounding quotes from values.
func WithoutStripQuotes() SanitizeOption {
	return func(s *Sanitizer) { s.stripQuotes = false }
}

// Sanitize processes a slice of Entry values and returns cleaned entries.
func (s *Sanitizer) Sanitize(entries []Entry) []Entry {
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		key := e.Key
		val := e.Value

		if s.trimKeys {
			key = strings.TrimSpace(key)
		}
		if s.normalizeKey {
			key = normalizeKeyString(key)
		}
		if s.trimValues {
			val = strings.TrimSpace(val)
		}
		if s.stripQuotes {
			val = stripSurroundingQuotes(val)
		}

		if key == "" {
			continue
		}
		out = append(out, Entry{Key: key, Value: val})
	}
	return out
}

func normalizeKeyString(k string) string {
	k = strings.ToUpper(k)
	var b strings.Builder
	for _, r := range k {
		if r == '-' || r == ' ' {
			b.WriteRune('_')
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func stripSurroundingQuotes(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') ||
			(v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}
