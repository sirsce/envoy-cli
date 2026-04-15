package env

import (
	"sort"
	"strings"
)

// ScoreRule assigns a numeric score to an entry based on custom logic.
type ScoreRule func(entry Entry) float64

// ScoredEntry pairs an entry with its computed score.
type ScoredEntry struct {
	Entry Entry
	Score float64
}

// Scorer computes and ranks env entries based on configurable rules.
type Scorer struct {
	rules []ScoreRule
}

// NewScorer creates a Scorer with the given scoring rules.
func NewScorer(rules ...ScoreRule) *Scorer {
	return &Scorer{rules: rules}
}

// Score computes a total score for each entry by summing all rule scores.
func (s *Scorer) Score(entries []Entry) []ScoredEntry {
	scored := make([]ScoredEntry, len(entries))
	for i, e := range entries {
		var total float64
		for _, rule := range s.rules {
			total += rule(e)
		}
		scored[i] = ScoredEntry{Entry: e, Score: total}
	}
	return scored
}

// Rank returns entries sorted by descending score.
func (s *Scorer) Rank(entries []Entry) []ScoredEntry {
	scored := s.Score(entries)
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	return scored
}

// HasEmptyValue is a built-in rule that penalises entries with empty values.
func HasEmptyValue(entry Entry) float64 {
	if strings.TrimSpace(entry.Value) == "" {
		return -1.0
	}
	return 0.0
}

// IsUppercaseKey is a built-in rule that rewards fully uppercased keys.
func IsUppercaseKey(entry Entry) float64 {
	if entry.Key == strings.ToUpper(entry.Key) {
		return 1.0
	}
	return 0.0
}

// HasPrefixRule returns a rule that rewards entries whose key starts with prefix.
func HasPrefixRule(prefix string) ScoreRule {
	return func(entry Entry) float64 {
		if strings.HasPrefix(entry.Key, prefix) {
			return 1.0
		}
		return 0.0
	}
}
