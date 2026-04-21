package env

import "strings"

// ClassifyResult holds the classification label and confidence for an entry.
type ClassifyResult struct {
	Key        string
	Label      string
	Confidence float64
}

// Classifier assigns semantic labels to env entries based on key patterns.
type Classifier struct {
	rules []classifyRule
}

type classifyRule struct {
	label    string
	matchers []string
}

// NewClassifier returns a Classifier with built-in rules.
func NewClassifier() *Classifier {
	return &Classifier{
		rules: []classifyRule{
			{label: "secret", matchers: []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "APIKEY", "PRIVATE_KEY"}},
			{label: "database", matchers: []string{"DB_", "DATABASE_", "POSTGRES", "MYSQL", "MONGO", "REDIS"}},
			{label: "network", matchers: []string{"HOST", "PORT", "URL", "URI", "ADDR", "ENDPOINT"}},
			{label: "feature_flag", matchers: []string{"FEATURE_", "FLAG_", "ENABLE_", "DISABLE_"}},
			{label: "observability", matchers: []string{"LOG_", "TRACE_", "METRIC_", "SENTRY_", "DATADOG_"}},
		},
	}
}

// Classify returns ClassifyResult for each entry.
func (c *Classifier) Classify(entries []Entry) []ClassifyResult {
	results := make([]ClassifyResult, 0, len(entries))
	for _, e := range entries {
		label, confidence := c.classify(e.Key)
		results = append(results, ClassifyResult{
			Key:        e.Key,
			Label:      label,
			Confidence: confidence,
		})
	}
	return results
}

func (c *Classifier) classify(key string) (string, float64) {
	upper := strings.ToUpper(key)
	for _, rule := range c.rules {
		for _, m := range rule.matchers {
			if strings.Contains(upper, m) {
				return rule.label, 0.9
			}
		}
	}
	return "generic", 0.5
}

// FilterByLabel returns only entries whose classification matches the given label.
func (c *Classifier) FilterByLabel(entries []Entry, label string) []Entry {
	results := c.Classify(entries)
	out := make([]Entry, 0)
	for i, r := range results {
		if r.Label == label {
			out = append(out, entries[i])
		}
	}
	return out
}
