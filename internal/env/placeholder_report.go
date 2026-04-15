package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// PlaceholderScanResult holds detected placeholders for a single entry.
type PlaceholderScanResult struct {
	Key          string   `json:"key"`
	Placeholders []string `json:"placeholders"`
}

// PlaceholderReport summarises placeholder detection across entries.
type PlaceholderReport struct {
	results []PlaceholderScanResult
}

// NewPlaceholderReport creates a report from a resolver and a set of entries.
func NewPlaceholderReport(r *PlaceholderResolver, entries []Entry) *PlaceholderReport {
	var results []PlaceholderScanResult
	for _, e := range entries {
		ph := r.DetectPlaceholders(e.Value)
		if len(ph) > 0 {
			results = append(results, PlaceholderScanResult{Key: e.Key, Placeholders: ph})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})
	return &PlaceholderReport{results: results}
}

// HasPlaceholders returns true if any placeholders were detected.
func (p *PlaceholderReport) HasPlaceholders() bool {
	return len(p.results) > 0
}

// WriteText writes a human-readable report to w.
func (p *PlaceholderReport) WriteText(w io.Writer) error {
	if len(p.results) == 0 {
		_, err := fmt.Fprintln(w, "No placeholders detected.")
		return err
	}
	for _, r := range p.results {
		if _, err := fmt.Fprintf(w, "%s: %v\n", r.Key, r.Placeholders); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes a JSON-encoded report to w.
func (p *PlaceholderReport) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(p.results)
}
