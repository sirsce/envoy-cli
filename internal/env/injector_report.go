package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// InjectorReport formats injection results for display.
type InjectorReport struct {
	results []InjectorResult
}

// NewInjectorReport creates a new InjectorReport from a slice of results.
func NewInjectorReport(results []InjectorResult) *InjectorReport {
	return &InjectorReport{results: results}
}

// WriteText writes a human-readable summary to w.
func (r *InjectorReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "No injection results.")
		return err
	}
	for _, res := range r.results {
		status := "INJECTED"
		if res.Skipped {
			status = "SKIPPED "
		}
		if _, err := fmt.Fprintf(w, "[%s] %s — %s\n", status, res.Key, res.Reason); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes results as a JSON array to w.
func (r *InjectorReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Key      string `json:"key"`
		Injected bool   `json:"injected"`
		Skipped  bool   `json:"skipped"`
		Reason   string `json:"reason"`
	}
	out := make([]jsonResult, len(r.results))
	for i, res := range r.results {
		out[i] = jsonResult{
			Key:      res.Key,
			Injected: res.Injected,
			Skipped:  res.Skipped,
			Reason:   res.Reason,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
