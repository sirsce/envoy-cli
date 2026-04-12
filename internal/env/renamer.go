package env

import "fmt"

// RenameRule describes a single key rename operation.
type RenameRule struct {
	From string
	To   string
}

// RenameResult records the outcome of a single rename operation.
type RenameResult struct {
	Rule    RenameRule
	Applied bool
	Reason  string
}

// Renamer applies a set of rename rules to a slice of entries.
type Renamer struct {
	rules []RenameRule
}

// NewRenamer creates a Renamer with the given rules.
func NewRenamer(rules []RenameRule) *Renamer {
	return &Renamer{rules: rules}
}

// Apply renames keys in entries according to the rules.
// It returns the updated entries and a result for each rule.
// If a target key already exists, the rule is skipped.
func (r *Renamer) Apply(entries []Entry) ([]Entry, []RenameResult, error) {
	results := make([]RenameResult, 0, len(r.rules))

	// Build a working map for fast lookup.
	index := make(map[string]int, len(entries))
	for i, e := range entries {
		index[e.Key] = i
	}

	for _, rule := range r.rules {
		if rule.From == "" || rule.To == "" {
			return nil, nil, fmt.Errorf("renamer: rule has empty From or To field")
		}

		srcIdx, srcExists := index[rule.From]
		if !srcExists {
			results = append(results, RenameResult{Rule: rule, Applied: false, Reason: "source key not found"})
			continue
		}

		if _, dstExists := index[rule.To]; dstExists {
			results = append(results, RenameResult{Rule: rule, Applied: false, Reason: "destination key already exists"})
			continue
		}

		// Perform the rename in-place.
		entries[srcIdx].Key = rule.To
		index[rule.To] = srcIdx
		delete(index, rule.From)

		results = append(results, RenameResult{Rule: rule, Applied: true})
	}

	return entries, results, nil
}

// AppliedCount returns the number of rules that were successfully applied.
func AppliedCount(results []RenameResult) int {
	n := 0
	for _, r := range results {
		if r.Applied {
			n++
		}
	}
	return n
}
