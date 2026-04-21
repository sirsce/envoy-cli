package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// StripReport formats StripResult slices as text or JSON.
type StripReport struct {
	results []StripResult
}

// NewStripReport creates a StripReport from a slice of StripResult.
func NewStripReport(results []StripResult) *StripReport {
	return &StripReport{results: results}
}

// WriteText writes a human-readable summary to w.
func (r *StripReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "No entries stripped.")
		return err
	}
	for _, res := range r.results {
		action := "removed"
		if !res.Stripped {
			continue
		}
		if res.OldVal == "" {
			action = "blanked"
		}
		if _, err := fmt.Fprintf(w, "  [stripped] %s (was: %q) — %s\n", res.Key, res.OldVal, action); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes a JSON array of results to w.
func (r *StripReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Key      string `json:"key"`
		OldValue string `json:"old_value"`
		Stripped bool   `json:"stripped"`
	}
	out := make([]jsonResult, len(r.results))
	for i, res := range r.results {
		out[i] = jsonResult{Key: res.Key, OldValue: res.OldVal, Stripped: res.Stripped}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
