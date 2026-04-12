package env

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaField describes a single expected env variable.
type SchemaField struct {
	Key         string
	Required    bool
	Default     string
	Pattern     string // optional regex pattern for value validation
	Description string
}

// Schema holds a collection of field definitions for an env file.
type Schema struct {
	fields map[string]SchemaField
}

// NewSchema creates a Schema from a slice of SchemaFields.
func NewSchema(fields []SchemaField) *Schema {
	s := &Schema{fields: make(map[string]SchemaField, len(fields))}
	for _, f := range fields {
		s.fields[f.Key] = f
	}
	return s
}

// SchemaViolation represents a single schema validation error.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks a map of env entries against the schema.
// It returns all violations found.
func (s *Schema) Validate(entries []Entry) []SchemaViolation {
	present := make(map[string]string, len(entries))
	for _, e := range entries {
		present[e.Key] = e.Value
	}

	var violations []SchemaViolation

	for key, field := range s.fields {
		val, exists := present[key]
		if !exists || strings.TrimSpace(val) == "" {
			if field.Required && field.Default == "" {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}
		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}
	}

	return violations
}

// ApplyDefaults returns a new slice of entries with defaults filled in for
// missing or empty keys defined in the schema.
func (s *Schema) ApplyDefaults(entries []Entry) []Entry {
	present := make(map[string]bool, len(entries))
	for _, e := range entries {
		if strings.TrimSpace(e.Value) != "" {
			present[e.Key] = true
		}
	}

	result := make([]Entry, len(entries))
	copy(result, entries)

	for key, field := range s.fields {
		if !present[key] && field.Default != "" {
			result = append(result, Entry{Key: key, Value: field.Default})
		}
	}
	return result
}
