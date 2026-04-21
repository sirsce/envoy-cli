package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// CastReport formats CastResult slices as text or JSON.
type CastReport struct {
	results []CastResult
}

// NewCastReport creates a CastReport from the given results.
func NewCastReport(results []CastResult) *CastReport {
	return &CastReport{results: results}
}

// WriteText writes a human-readable summary to w.
func (r *CastReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "No cast operations performed.")
		return err
	}
	for _, res := range r.results {
		if res.Skipped {
			if _, err := fmt.Fprintf(w, "SKIP  %-24s  error: %v\n", res.Key, res.Err); err != nil {
				return err
			}
			continue
		}
		if _, err := fmt.Fprintf(w, "CAST  %-24s  %s → %s  (type: %s)\n",
			res.Key, res.From, res.To, res.Type); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes a JSON array of results to w.
func (r *CastReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Key     string `json:"key"`
		From    string `json:"from"`
		To      string `json:"to,omitempty"`
		Type    string `json:"type"`
		Skipped bool   `json:"skipped"`
		Error   string `json:"error,omitempty"`
	}

	payload := make([]jsonResult, len(r.results))
	for i, res := range r.results {
		je := jsonResult{
			Key:     res.Key,
			From:    res.From,
			To:      res.To,
			Type:    string(res.Type),
			Skipped: res.Skipped,
		}
		if res.Err != nil {
			je.Error = res.Err.Error()
		}
		payload[i] = je
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
