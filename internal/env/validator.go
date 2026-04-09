package env

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a single validation issue.
type ValidationError struct {
	Line    int
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: key %q: %s", e.Line, e.Key, e.Message)
	}
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

func (r *ValidationResult) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "; ")
}

var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Validator checks .env entries for common issues.
type Validator struct {
	requiredKeys []string
}

// NewValidator creates a Validator. Pass required key names to enforce presence.
func NewValidator(requiredKeys ...string) *Validator {
	return &Validator{requiredKeys: requiredKeys}
}

// Validate checks a slice of Entry values and returns a ValidationResult.
func (v *Validator) Validate(entries []Entry) ValidationResult {
	result := ValidationResult{}
	seen := make(map[string]int) // key -> first line seen

	for i, entry := range entries {
		lineNum := i + 1
		if !validKeyPattern.MatchString(entry.Key) {
			result.Errors = append(result.Errors, ValidationError{
				Line:    lineNum,
				Key:     entry.Key,
				Message: "invalid key name (must match [A-Za-z_][A-Za-z0-9_]*)",
			})
		}
		if prev, dup := seen[entry.Key]; dup {
			result.Errors = append(result.Errors, ValidationError{
				Line:    lineNum,
				Key:     entry.Key,
				Message: fmt.Sprintf("duplicate key (first defined on line %d)", prev),
			})
		} else {
			seen[entry.Key] = lineNum
		}
	}

	for _, req := range v.requiredKeys {
		if _, ok := seen[req]; !ok {
			result.Errors = append(result.Errors, ValidationError{
				Key:     req,
				Message: "required key is missing",
			})
		}
	}

	return result
}
