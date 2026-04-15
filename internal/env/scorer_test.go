package env

import (
	"testing"
)

func makeScorerEntries() []Entry {
	return []Entry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
		{Key: "api_key", Value: ""},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	}
}

func TestScorer_ScoreAllZero(t *testing.T) {
	scorer := NewScorer()
	entries := makeScorerEntries()
	scored := scorer.Score(entries)
	if len(scored) != len(entries) {
		t.Fatalf("expected %d scored entries, got %d", len(entries), len(scored))
	}
	for _, s := range scored {
		if s.Score != 0 {
			t.Errorf("expected score 0 for %s, got %f", s.Entry.Key, s.Score)
		}
	}
}

func TestScorer_UppercaseKeyRule(t *testing.T) {
	scorer := NewScorer(IsUppercaseKey)
	entries := makeScorerEntries()
	scored := scorer.Score(entries)

	scoreMap := map[string]float64{}
	for _, s := range scored {
		scoreMap[s.Entry.Key] = s.Score
	}

	if scoreMap["DATABASE_URL"] != 1.0 {
		t.Errorf("DATABASE_URL should score 1.0, got %f", scoreMap["DATABASE_URL"])
	}
	if scoreMap["api_key"] != 0.0 {
		t.Errorf("api_key should score 0.0, got %f", scoreMap["api_key"])
	}
}

func TestScorer_EmptyValuePenalty(t *testing.T) {
	scorer := NewScorer(HasEmptyValue)
	entries := makeScorerEntries()
	scored := scorer.Score(entries)

	for _, s := range scored {
		if s.Entry.Key == "api_key" && s.Score != -1.0 {
			t.Errorf("api_key should score -1.0, got %f", s.Score)
		}
		if s.Entry.Key == "DATABASE_URL" && s.Score != 0.0 {
			t.Errorf("DATABASE_URL should score 0.0, got %f", s.Score)
		}
	}
}

func TestScorer_Rank_Descending(t *testing.T) {
	scorer := NewScorer(IsUppercaseKey, HasEmptyValue)
	entries := makeScorerEntries()
	ranked := scorer.Rank(entries)

	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score > ranked[i-1].Score {
			t.Errorf("rank not descending at index %d: %f > %f", i, ranked[i].Score, ranked[i-1].Score)
		}
	}
}

func TestScorer_HasPrefixRule(t *testing.T) {
	scorer := NewScorer(HasPrefixRule("DB_"))
	entries := makeScorerEntries()
	scored := scorer.Score(entries)

	scoreMap := map[string]float64{}
	for _, s := range scored {
		scoreMap[s.Entry.Key] = s.Score
	}

	if scoreMap["DB_HOST"] != 1.0 {
		t.Errorf("DB_HOST should score 1.0, got %f", scoreMap["DB_HOST"])
	}
	if scoreMap["DATABASE_URL"] != 0.0 {
		t.Errorf("DATABASE_URL should score 0.0, got %f", scoreMap["DATABASE_URL"])
	}
}

func TestScorer_MultipleRules_Combined(t *testing.T) {
	scorer := NewScorer(IsUppercaseKey, HasEmptyValue, HasPrefixRule("DB_"))
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
	}
	scored := scorer.Score(entries)
	// IsUppercaseKey: +1, HasEmptyValue: 0, HasPrefixRule(DB_): +1 => 2.0
	if scored[0].Score != 2.0 {
		t.Errorf("expected combined score 2.0, got %f", scored[0].Score)
	}
}
