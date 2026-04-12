package env

import (
	"strings"
)

// MaskStyle controls how values are masked.
type MaskStyle int

const (
	// MaskFull replaces the entire value with asterisks.
	MaskFull MaskStyle = iota
	// MaskPartial reveals the first and last characters.
	MaskPartial
	// MaskHash replaces the value with a fixed placeholder.
	MaskHash
)

// MaskerOption configures a Masker.
type MaskerOption func(*Masker)

// Masker masks sensitive entry values based on key patterns.
type Masker struct {
	keys  map[string]struct{}
	style MaskStyle
}

// NewMasker returns a Masker with the given style and sensitive key names.
func NewMasker(style MaskStyle, sensitiveKeys ...string) *Masker {
	m := &Masker{
		keys:  make(map[string]struct{}, len(sensitiveKeys)),
		style: style,
	}
	for _, k := range sensitiveKeys {
		m.keys[strings.ToUpper(k)] = struct{}{}
	}
	return m
}

// IsSensitive reports whether the given key is considered sensitive.
func (m *Masker) IsSensitive(key string) bool {
	_, ok := m.keys[strings.ToUpper(key)]
	return ok
}

// Mask returns a masked version of value if the key is sensitive,
// otherwise it returns the value unchanged.
func (m *Masker) Mask(key, value string) string {
	if !m.IsSensitive(key) {
		return value
	}
	if value == "" {
		return ""
	}
	switch m.style {
	case MaskPartial:
		if len(value) <= 2 {
			return strings.Repeat("*", len(value))
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	case MaskHash:
		return "[REDACTED]"
	default: // MaskFull
		return strings.Repeat("*", len(value))
	}
}

// Apply returns a new slice of Entry with sensitive values masked.
func (m *Masker) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		out[i] = Entry{
			Key:   e.Key,
			Value: m.Mask(e.Key, e.Value),
		}
	}
	return out
}
