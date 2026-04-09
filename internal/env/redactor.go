package env

import (
	"regexp"
	"strings"
)

// RedactMode controls how sensitive values are displayed.
type RedactMode int

const (
	// RedactFull replaces the entire value with asterisks.
	RedactFull RedactMode = iota
	// RedactPartial shows the first and last character with asterisks in between.
	RedactPartial
)

// sensitivePatterns holds regex patterns for keys considered sensitive.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|secret|token|key|api_key|auth|credential|private)`),
}

// Redactor masks sensitive environment variable values.
type Redactor struct {
	mode     RedactMode
	extraKeys []string
}

// NewRedactor creates a new Redactor with the given mode.
func NewRedactor(mode RedactMode, extraKeys ...string) *Redactor {
	return &Redactor{
		mode:      mode,
		extraKeys: extraKeys,
	}
}

// IsSensitive reports whether the given key name is considered sensitive.
func (r *Redactor) IsSensitive(key string) bool {
	for _, pattern := range sensitivePatterns {
		if pattern.MatchString(key) {
			return true
		}
	}
	for _, extra := range r.extraKeys {
		if strings.EqualFold(key, extra) {
			return true
		}
	}
	return false
}

// Redact returns a masked version of value if the key is sensitive.
// If the key is not sensitive, the original value is returned unchanged.
func (r *Redactor) Redact(key, value string) string {
	if !r.IsSensitive(key) {
		return value
	}
	if len(value) == 0 {
		return value
	}
	switch r.mode {
	case RedactPartial:
		if len(value) <= 2 {
			return strings.Repeat("*", len(value))
		}
		return string(value[0]) + strings.Repeat("*", len(value)-2) + string(value[len(value)-1])
	default:
		return strings.Repeat("*", len(value))
	}
}

// RedactEntries returns a copy of entries with sensitive values masked.
func (r *Redactor) RedactEntries(entries []Entry) []Entry {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		result[i] = Entry{
			Key:   e.Key,
			Value: r.Redact(e.Key, e.Value),
		}
	}
	return result
}
