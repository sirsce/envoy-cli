package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ClassifierReport renders classification results.
type ClassifierReport struct {
	results []ClassifyResult
}

// NewClassifierReport creates a report from classification results.
func NewClassifierReport(results []ClassifyResult) *ClassifierReport {
	return &ClassifierReport{results: results}
}

// WriteText writes a human-readable classification summary to w.
func (r *ClassifierReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "no entries classified")
		return err
	}
	sorted := make([]ClassifyResult, len(r.results))
	copy(sorted, r.results)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Label == sorted[j].Label {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Label < sorted[j].Label
	})
	for _, res := range sorted {
		_, err := fmt.Fprintf(w, "%-40s %-16s (%.1f)\n", res.Key, res.Label, res.Confidence)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes classification results as JSON to w.
func (r *ClassifierReport) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.results)
}

// Summary returns a map of label -> count.
func (r *ClassifierReport) Summary() map[string]int {
	m := make(map[string]int)
	for _, res := range r.results {
		m[res.Label]++
	}
	return m
}
