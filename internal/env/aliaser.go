package env

import (
	"errors"
	"fmt"
)

// AliasRule maps an alias key to a source key.
type AliasRule struct {
	Alias  string
	Source string
}

// AliasResult records the outcome of applying an alias.
type AliasResult struct {
	Alias   string
	Source  string
	Applied bool
	Reason  string
}

// Aliaser creates virtual alias entries that mirror the value of a source key.
type Aliaser struct {
	rules     []AliasRule
	overwrite bool
}

// NewAliaser constructs an Aliaser with the given rules.
// If overwrite is true, existing entries with the alias key will be replaced.
func NewAliaser(rules []AliasRule, overwrite bool) (*Aliaser, error) {
	for _, r := range rules {
		if r.Alias == "" || r.Source == "" {
			return nil, errors.New("aliaser: alias and source keys must not be empty")
		}
		if r.Alias == r.Source {
			return nil, fmt.Errorf("aliaser: alias %q must differ from source", r.Alias)
		}
	}
	return &Aliaser{rules: rules, overwrite: overwrite}, nil
}

// Apply processes entries and returns a new slice with alias entries appended or replaced.
func (a *Aliaser) Apply(entries []Entry) ([]Entry, []AliasResult) {
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	out := make([]Entry, len(entries))
	copy(out, entries)

	var results []AliasResult

	for _, rule := range a.rules {
		srcIdx, srcFound := index[rule.Source]
		if !srcFound {
			results = append(results, AliasResult{
				Alias:   rule.Alias,
				Source:  rule.Source,
				Applied: false,
				Reason:  "source key not found",
			})
			continue
		}

		if existIdx, exists := index[rule.Alias]; exists {
			if !a.overwrite {
				results = append(results, AliasResult{
					Alias:   rule.Alias,
					Source:  rule.Source,
					Applied: false,
					Reason:  "alias key already exists",
				})
				continue
			}
			out[existIdx] = Entry{Key: rule.Alias, Value: out[srcIdx].Value}
			results = append(results, AliasResult{Alias: rule.Alias, Source: rule.Source, Applied: true, Reason: "overwritten"})
			continue
		}

		newEntry := Entry{Key: rule.Alias, Value: out[srcIdx].Value}
		index[rule.Alias] = len(out)
		out = append(out, newEntry)
		results = append(results, AliasResult{Alias: rule.Alias, Source: rule.Source, Applied: true, Reason: "created"})
	}

	return out, results
}

// AppliedCount returns the number of successfully applied alias results.
func AppliedAliasCount(results []AliasResult) int {
	n := 0
	for _, r := range results {
		if r.Applied {
			n++
		}
	}
	return n
}
