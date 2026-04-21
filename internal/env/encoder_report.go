package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// EncodeReport formats encoding results as text or JSON.
type EncodeReport struct {
	results []EncodeResult
}

// NewEncodeReport creates a report from a slice of EncodeResult.
func NewEncodeReport(results []EncodeResult) *EncodeReport {
	return &EncodeReport{results: results}
}

// WriteText writes a human-readable summary to w.
func (r *EncodeReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "No entries encoded.")
		return err
	}
	for _, res := range r.results {
		if res.Skipped {
			fmt.Fprintf(w, "  SKIP  %-30s\n", res.Key)
		} else {
			fmt.Fprintf(w, "  OK    %-30s %s -> %s\n", res.Key, res.Original, res.Encoded)
		}
	}
	return nil
}

// WriteJSON writes results as a JSON array to w.
func (r *EncodeReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Key      string `json:"key"`
		Original string `json:"original"`
		Encoded  string `json:"encoded"`
		Skipped  bool   `json:"skipped"`
	}
	var out []jsonResult
	for _, res := range r.results {
		out = append(out, jsonResult{
			Key:      res.Key,
			Original: res.Original,
			Encoded:  res.Encoded,
			Skipped:  res.Skipped,
		})
	}
	if out == nil {
		out = []jsonResult{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
