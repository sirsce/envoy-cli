package env

import (
	"testing"
)

func makeClassifierEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "FEATURE_DARK_MODE", Value: "true"},
		{Key: "LOG_LEVEL", Value: "info"},
		{Key: "APP_NAME", Value: "envoy"},
	}
}

func TestClassifier_LabelSecret(t *testing.T) {
	c := NewClassifier()
	entries := makeClassifierEntries()
	results := c.Classify(entries)
	for _, r := range results {
		if r.Key == "API_KEY" && r.Label != "secret" {
			t.Errorf("expected API_KEY to be secret, got %s", r.Label)
		}
		if r.Key == "DB_PASSWORD" && r.Label != "secret" {
			t.Errorf("expected DB_PASSWORD to be secret, got %s", r.Label)
		}
	}
}

func TestClassifier_LabelDatabase(t *testing.T) {
	c := NewClassifier()
	results := c.Classify([]Entry{{Key: "DB_HOST", Value: "localhost"}})
	if len(results) != 1 || results[0].Label != "database" {
		t.Errorf("expected database label, got %s", results[0].Label)
	}
}

func TestClassifier_LabelFeatureFlag(t *testing.T) {
	c := NewClassifier()
	results := c.Classify([]Entry{{Key: "FEATURE_DARK_MODE", Value: "true"}})
	if results[0].Label != "feature_flag" {
		t.Errorf("expected feature_flag, got %s", results[0].Label)
	}
}

func TestClassifier_LabelGeneric(t *testing.T) {
	c := NewClassifier()
	results := c.Classify([]Entry{{Key: "APP_NAME", Value: "envoy"}})
	if results[0].Label != "generic" {
		t.Errorf("expected generic, got %s", results[0].Label)
	}
	if results[0].Confidence != 0.5 {
		t.Errorf("expected confidence 0.5, got %f", results[0].Confidence)
	}
}

func TestClassifier_FilterByLabel(t *testing.T) {
	c := NewClassifier()
	entries := makeClassifierEntries()
	secrets := c.FilterByLabel(entries, "secret")
	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}
}

func TestClassifier_AllEntriesClassified(t *testing.T) {
	c := NewClassifier()
	entries := makeClassifierEntries()
	results := c.Classify(entries)
	if len(results) != len(entries) {
		t.Errorf("expected %d results, got %d", len(entries), len(results))
	}
	for _, r := range results {
		if r.Label == "" {
			t.Errorf("entry %s has empty label", r.Key)
		}
	}
}
