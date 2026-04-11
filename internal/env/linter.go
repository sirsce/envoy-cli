package env

import (
	"fmt"
	"strings"
)

// LintRule represents a single linting rule.
type LintRule struct {
	Name    string
	Message string
	Check   func(entry Entry) bool
}

// LintResult holds the outcome of a lint check for a single entry.
type LintResult struct {
	Entry   Entry
	Rule    string
	Message string
}

// Linter checks env entries against a set of rules.
type Linter struct {
	rules []LintRule
}

// NewLinter creates a Linter with the default built-in rules.
func NewLinter() *Linter {
	return &Linter{
		rules: defaultRules(),
	}
}

// AddRule appends a custom rule to the linter.
func (l *Linter) AddRule(rule LintRule) {
	l.rules = append(l.rules, rule)
}

// Lint runs all rules against the given entries and returns any violations.
func (l *Linter) Lint(entries []Entry) []LintResult {
	var results []LintResult
	for _, entry := range entries {
		for _, rule := range l.rules {
			if !rule.Check(entry) {
				results = append(results, LintResult{
					Entry:   entry,
					Rule:    rule.Name,
					Message: fmt.Sprintf("%s: key=%q", rule.Message, entry.Key),
				})
			}
		}
	}
	return results
}

// HasViolations returns true if any lint violations were found.
func HasViolations(results []LintResult) bool {
	return len(results) > 0
}

func defaultRules() []LintRule {
	return []LintRule{
		{
			Name:    "no-empty-value",
			Message: "value should not be empty",
			Check: func(e Entry) bool {
				return strings.TrimSpace(e.Value) != ""
			},
		},
		{
			Name:    "uppercase-key",
			Message: "key should be uppercase",
			Check: func(e Entry) bool {
				return e.Key == strings.ToUpper(e.Key)
			},
		},
		{
			Name:    "no-whitespace-in-key",
			Message: "key must not contain whitespace",
			Check: func(e Entry) bool {
				return !strings.ContainsAny(e.Key, " \t")
			},
		},
	}
}
