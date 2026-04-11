package env

import "strings"

// TransformFunc is a function that transforms an entry's value.
type TransformFunc func(value string) string

// Transformer applies a chain of transformations to env entries.
type Transformer struct {
	transforms []TransformFunc
}

// NewTransformer creates a new Transformer with the given transform functions.
func NewTransformer(fns ...TransformFunc) *Transformer {
	return &Transformer{transforms: fns}
}

// Apply applies all transformations to the given entries and returns new entries.
func (t *Transformer) Apply(entries []Entry) []Entry {
	result := make([]Entry, len(entries))
	for i, e := range entries {
		val := e.Value
		for _, fn := range t.transforms {
			val = fn(val)
		}
		result[i] = Entry{Key: e.Key, Value: val}
	}
	return result
}

// TrimSpaceTransform removes leading and trailing whitespace from values.
func TrimSpaceTransform() TransformFunc {
	return func(value string) string {
		return strings.TrimSpace(value)
	}
}

// ToUpperTransform converts values to uppercase.
func ToUpperTransform() TransformFunc {
	return func(value string) string {
		return strings.ToUpper(value)
	}
}

// ToLowerTransform converts values to lowercase.
func ToLowerTransform() TransformFunc {
	return func(value string) string {
		return strings.ToLower(value)
	}
}

// ReplaceTransform replaces all occurrences of old with new in values.
func ReplaceTransform(old, new string) TransformFunc {
	return func(value string) string {
		return strings.ReplaceAll(value, old, new)
	}
}

// PrefixValueTransform prepends a prefix to each value.
func PrefixValueTransform(prefix string) TransformFunc {
	return func(value string) string {
		if value == "" {
			return value
		}
		return prefix + value
	}
}
