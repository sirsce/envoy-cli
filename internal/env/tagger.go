package env

import (
	"fmt"
	"sort"
	"strings"
)

// Tag represents a label attached to an env entry.
type Tag struct {
	Name  string
	Value string
}

// Tagger manages tags associated with env entry keys.
type Tagger struct {
	tags map[string][]Tag // key -> tags
}

// NewTagger creates a new Tagger instance.
func NewTagger() *Tagger {
	return &Tagger{tags: make(map[string][]Tag)}
}

// Tag attaches a named tag (with optional value) to an env key.
func (t *Tagger) Tag(key, name, value string) {
	key = strings.TrimSpace(key)
	name = strings.TrimSpace(name)
	if key == "" || name == "" {
		return
	}
	t.tags[key] = append(t.tags[key], Tag{Name: name, Value: value})
}

// GetTags returns all tags for a given key.
func (t *Tagger) GetTags(key string) []Tag {
	return t.tags[key]
}

// HasTag reports whether a key has a tag with the given name.
func (t *Tagger) HasTag(key, name string) bool {
	for _, tag := range t.tags[key] {
		if tag.Name == name {
			return true
		}
	}
	return false
}

// FilterByTag returns entries whose keys carry the specified tag name.
func (t *Tagger) FilterByTag(entries []Entry, name string) []Entry {
	var out []Entry
	for _, e := range entries {
		if t.HasTag(e.Key, name) {
			out = append(out, e)
		}
	}
	return out
}

// RemoveTag removes all tags with the given name from a key.
func (t *Tagger) RemoveTag(key, name string) {
	existing := t.tags[key]
	filtered := existing[:0]
	for _, tag := range existing {
		if tag.Name != name {
			filtered = append(filtered, tag)
		}
	}
	t.tags[key] = filtered
}

// Keys returns all keys that have at least one tag, sorted.
func (t *Tagger) Keys() []string {
	keys := make([]string, 0, len(t.tags))
	for k := range t.tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Summary returns a human-readable summary of tags for a key.
func (t *Tagger) Summary(key string) string {
	tags := t.tags[key]
	if len(tags) == 0 {
		return fmt.Sprintf("%s: (no tags)", key)
	}
	parts := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag.Value != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", tag.Name, tag.Value))
		} else {
			parts = append(parts, tag.Name)
		}
	}
	return fmt.Sprintf("%s: [%s]", key, strings.Join(parts, ", "))
}
