package env

import (
	"fmt"
	"regexp"
	"strings"
)

// PlaceholderStyle defines how placeholders are detected.
type PlaceholderStyle int

const (
	StyleDoubleBrace  PlaceholderStyle = iota // {{ KEY }}
	StyleAngleBracket                         // <KEY>
	StylePercent                              // %KEY%
)

var (
	reDoubleBrace  = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)
	reAngleBracket = regexp.MustCompile(`<(\w+)>`)
	rePercent      = regexp.MustCompile(`%(\w+)%`)
)

// PlaceholderResolver replaces placeholder tokens in entry values
// using a provided lookup map.
type PlaceholderResolver struct {
	style   PlaceholderStyle
	strict  bool
	pattern *regexp.Regexp
}

// NewPlaceholderResolver creates a resolver for the given style.
// If strict is true, unresolved placeholders return an error.
func NewPlaceholderResolver(style PlaceholderStyle, strict bool) *PlaceholderResolver {
	var pat *regexp.Regexp
	switch style {
	case StyleAngleBracket:
		pat = reAngleBracket
	case StylePercent:
		pat = rePercent
	default:
		pat = reDoubleBrace
	}
	return &PlaceholderResolver{style: style, strict: strict, pattern: pat}
}

// Resolve replaces placeholders in all entry values using the provided lookup.
func (r *PlaceholderResolver) Resolve(entries []Entry, lookup map[string]string) ([]Entry, error) {
	out := make([]Entry, len(entries))
	for i, e := range entries {
		resolved, err := r.resolveValue(e.Value, lookup)
		if err != nil {
			return nil, fmt.Errorf("entry %q: %w", e.Key, err)
		}
		out[i] = Entry{Key: e.Key, Value: resolved}
	}
	return out, nil
}

// DetectPlaceholders returns the list of placeholder keys found in a value.
func (r *PlaceholderResolver) DetectPlaceholders(value string) []string {
	matches := r.pattern.FindAllStringSubmatch(value, -1)
	keys := make([]string, 0, len(matches))
	for _, m := range matches {
		keys = append(keys, m[1])
	}
	return keys
}

func (r *PlaceholderResolver) resolveValue(value string, lookup map[string]string) (string, error) {
	var resolveErr error
	result := r.pattern.ReplaceAllStringFunc(value, func(match string) string {
		sub := r.pattern.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		key := strings.TrimSpace(sub[1])
		if val, ok := lookup[key]; ok {
			return val
		}
		if r.strict {
			resolveErr = fmt.Errorf("unresolved placeholder %q", key)
			return match
		}
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}
