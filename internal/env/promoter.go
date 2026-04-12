package env

import "fmt"

// PromoteRule defines a promotion from one environment scope to another.
type PromoteRule struct {
	From string
	To   string
}

// PromoteResult records the outcome of a single key promotion.
type PromoteResult struct {
	Key       string
	From      string
	To        string
	Overwrote bool
	Skipped   bool
	Reason    string
}

// Promoter copies env entries between named scopes.
type Promoter struct {
	scopes    map[string][]Entry
	overwrite bool
}

// NewPromoter creates a Promoter with the given named scopes.
func NewPromoter(scopes map[string][]Entry, overwrite bool) *Promoter {
	return &Promoter{scopes: scopes, overwrite: overwrite}
}

// Promote copies entries from the source scope to the destination scope
// according to the provided rule, returning per-key results.
func (p *Promoter) Promote(rule PromoteRule) ([]PromoteResult, error) {
	src, ok := p.scopes[rule.From]
	if !ok {
		return nil, fmt.Errorf("promoter: source scope %q not found", rule.From)
	}
	if _, ok := p.scopes[rule.To]; !ok {
		p.scopes[rule.To] = []Entry{}
	}

	destMap := make(map[string]int)
	for i, e := range p.scopes[rule.To] {
		destMap[e.Key] = i
	}

	var results []PromoteResult
	dest := p.scopes[rule.To]

	for _, entry := range src {
		result := PromoteResult{Key: entry.Key, From: rule.From, To: rule.To}
		if idx, exists := destMap[entry.Key]; exists {
			if p.overwrite {
				dest[idx].Value = entry.Value
				result.Overwrote = true
			} else {
				result.Skipped = true
				result.Reason = "key already exists in destination"
			}
		} else {
			dest = append(dest, entry)
			destMap[entry.Key] = len(dest) - 1
		}
		results = append(results, result)
	}

	p.scopes[rule.To] = dest
	return results, nil
}

// Scope returns a copy of the entries for the named scope.
func (p *Promoter) Scope(name string) ([]Entry, bool) {
	entries, ok := p.scopes[name]
	if !ok {
		return nil, false
	}
	out := make([]Entry, len(entries))
	copy(out, entries)
	return out, true
}
