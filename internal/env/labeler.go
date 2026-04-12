package env

import (
	"fmt"
	"sort"
)

// Label represents a key-value metadata tag attached to an env entry.
type Label struct {
	Key   string
	Value string
}

// LabelSet holds labels for a given env entry key.
type LabelSet map[string]string

// Labeler manages arbitrary metadata labels on env entry keys.
type Labeler struct {
	labels map[string]LabelSet
}

// NewLabeler creates a new Labeler instance.
func NewLabeler() *Labeler {
	return &Labeler{
		labels: make(map[string]LabelSet),
	}
}

// Set attaches a label key=value to the given env entry key.
func (l *Labeler) Set(entryKey, labelKey, labelValue string) error {
	if entryKey == "" {
		return fmt.Errorf("labeler: entry key must not be empty")
	}
	if labelKey == "" {
		return fmt.Errorf("labeler: label key must not be empty")
	}
	if _, ok := l.labels[entryKey]; !ok {
		l.labels[entryKey] = make(LabelSet)
	}
	l.labels[entryKey][labelKey] = labelValue
	return nil
}

// Get returns all labels for a given env entry key.
func (l *Labeler) Get(entryKey string) LabelSet {
	if ls, ok := l.labels[entryKey]; ok {
		copy := make(LabelSet, len(ls))
		for k, v := range ls {
			copy[k] = v
		}
		return copy
	}
	return LabelSet{}
}

// Remove deletes a specific label from an env entry key.
func (l *Labeler) Remove(entryKey, labelKey string) {
	if ls, ok := l.labels[entryKey]; ok {
		delete(ls, labelKey)
		if len(ls) == 0 {
			delete(l.labels, entryKey)
		}
	}
}

// FilterByLabel returns entry keys that have a matching label key=value.
func (l *Labeler) FilterByLabel(labelKey, labelValue string) []string {
	var result []string
	for entryKey, ls := range l.labels {
		if v, ok := ls[labelKey]; ok && v == labelValue {
			result = append(result, entryKey)
		}
	}
	sort.Strings(result)
	return result
}

// Keys returns all env entry keys that have at least one label.
func (l *Labeler) Keys() []string {
	keys := make([]string, 0, len(l.labels))
	for k := range l.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
