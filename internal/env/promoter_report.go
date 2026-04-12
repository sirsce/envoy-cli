package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// PromoteReport formats promotion results for display.
type PromoteReport struct {
	results []PromoteResult
}

// NewPromoteReport creates a PromoteReport from a slice of results.
func NewPromoteReport(results []PromoteResult) *PromoteReport {
	return &PromoteReport{results: results}
}

// WriteText writes a human-readable summary to w.
func (r *PromoteReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "No keys promoted.")
		return err
	}
	for _, res := range r.results {
		switch {
		case res.Skipped:
			fmt.Fprintf(w, "SKIP  %-30s %s -> %s (%s)\n", res.Key, res.From, res.To, res.Reason)
		case res.Overwrote:
			fmt.Fprintf(w, "OVER  %-30s %s -> %s\n", res.Key, res.From, res.To)
		default:
			fmt.Fprintf(w, "OK    %-30s %s -> %s\n", res.Key, res.From, res.To)
		}
	}
	return nil
}

// WriteJSON writes results as a JSON array to w.
func (r *PromoteReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Key       string `json:"key"`
		From      string `json:"from"`
		To        string `json:"to"`
		Overwrote bool   `json:"overwrote,omitempty"`
		Skipped   bool   `json:"skipped,omitempty"`
		Reason    string `json:"reason,omitempty"`
	}
	out := make([]jsonResult, len(r.results))
	for i, res := range r.results {
		out[i] = jsonResult{
			Key:       res.Key,
			From:      res.From,
			To:        res.To,
			Overwrote: res.Overwrote,
			Skipped:   res.Skipped,
			Reason:    res.Reason,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// Summary returns counts of promoted, skipped, and overwritten keys.
func (r *PromoteReport) Summary() (promoted, skipped, overwrote int) {
	for _, res := range r.results {
		switch {
		case res.Skipped:
			skipped++
		case res.Overwrote:
			overwrote++
		default:
			promoted++
		}
	}
	return
}
