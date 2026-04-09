package env

import (
	"bytes"
	"fmt"
	"strings"
)

// Template represents a .env template with required and optional keys.
type Template struct {
	entries  []Entry
	required map[string]bool
	defaults map[string]string
}

// NewTemplate creates a new Template from a slice of entries.
func NewTemplate(entries []Entry) *Template {
	t := &Template{
		entries:  entries,
		required: make(map[string]bool),
		defaults: make(map[string]string),
	}
	for _, e := range entries {
		if e.Value == "" {
			t.required[e.Key] = true
		} else {
			t.defaults[e.Key] = e.Value
		}
	}
	return t
}

// Apply fills in the template with the provided values map.
// Returns a new slice of entries with values applied.
// Missing required keys are returned as an error.
func (t *Template) Apply(values map[string]string) ([]Entry, error) {
	var missing []string
	result := make([]Entry, 0, len(t.entries))

	for _, e := range t.entries {
		entry := Entry{Key: e.Key, Comment: e.Comment}
		if v, ok := values[e.Key]; ok {
			entry.Value = v
		} else if def, ok := t.defaults[e.Key]; ok {
			entry.Value = def
		} else if t.required[e.Key] {
			missing = append(missing, e.Key)
			continue
		}
		result = append(result, entry)
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required keys: %s", strings.Join(missing, ", "))
	}
	return result, nil
}

// Render returns a string representation of the template with placeholder comments.
func (t *Template) Render() string {
	var buf bytes.Buffer
	for _, e := range t.entries {
		if e.Comment != "" {
			fmt.Fprintf(&buf, "# %s\n", e.Comment)
		}
		if t.required[e.Key] {
			fmt.Fprintf(&buf, "%s=\n", e.Key)
		} else {
			fmt.Fprintf(&buf, "%s=%s\n", e.Key, e.Value)
		}
	}
	return buf.String()
}

// RequiredKeys returns the list of keys that must be provided.
func (t *Template) RequiredKeys() []string {
	keys := make([]string, 0, len(t.required))
	for k := range t.required {
		keys = append(keys, k)
	}
	return keys
}
