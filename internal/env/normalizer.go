package env

import (
	"strings"
)

// NormalizeMode controls how keys are normalized.
type NormalizeMode int

const (
	NormalizeModeUpper NormalizeMode = iota
	NormalizeModeSnake
	NormalizeModeCamel
)

// NormalizerOption configures a Normalizer.
type NormalizerOption func(*Normalizer)

// Normalizer applies key normalization rules to a slice of entries.
type Normalizer struct {
	mode        NormalizeMode
	modified    int
}

// WithNormalizeMode sets the normalization mode.
func WithNormalizeMode(m NormalizeMode) NormalizerOption {
	return func(n *Normalizer) {
		n.mode = m
	}
}

// NewNormalizer creates a Normalizer with the given options.
func NewNormalizer(opts ...NormalizerOption) *Normalizer {
	n := &Normalizer{mode: NormalizeModeUpper}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Normalize applies the normalization mode to all entry keys.
// It returns a new slice of entries; the originals are not mutated.
func (n *Normalizer) Normalize(entries []Entry) []Entry {
	n.modified = 0
	out := make([]Entry, len(entries))
	for i, e := range entries {
		newKey := n.applyMode(e.Key)
		if newKey != e.Key {
			n.modified++
		}
		out[i] = Entry{Key: newKey, Value: e.Value}
	}
	return out
}

// ModifiedCount returns the number of keys changed in the last Normalize call.
func (n *Normalizer) ModifiedCount() int {
	return n.modified
}

func (n *Normalizer) applyMode(key string) string {
	switch n.mode {
	case NormalizeModeUpper:
		return strings.ToUpper(key)
	case NormalizeModeSnake:
		return toSnakeCase(key)
	case NormalizeModeCamel:
		return toCamelCase(key)
	default:
		return key
	}
}

func toSnakeCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

func toCamelCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, p := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(p)
		} else if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		}
	}
	return strings.Join(parts, "")
}
