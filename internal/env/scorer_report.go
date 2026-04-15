package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// ScoreReport renders scored entry results in text or JSON format.
type ScoreReport struct {
	scored []ScoredEntry
}

// NewScoreReport creates a ScoreReport from a slice of ScoredEntry.
func NewScoreReport(scored []ScoredEntry) *ScoreReport {
	return &ScoreReport{scored: scored}
}

// WriteText writes a human-readable table of scored entries to w.
func (r *ScoreReport) WriteText(w io.Writer) error {
	if len(r.scored) == 0 {
		_, err := fmt.Fprintln(w, "No scored entries.")
		return err
	}
	_, err := fmt.Fprintf(w, "%-40s %8s\n", "KEY", "SCORE")
	if err != nil {
		return err
	}
	for _, s := range r.scored {
		_, err = fmt.Fprintf(w, "%-40s %8.2f\n", s.Entry.Key, s.Score)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes scored entries as a JSON array to w.
func (r *ScoreReport) WriteJSON(w io.Writer) error {
	type jsonEntry struct {
		Key   string  `json:"key"`
		Value string  `json:"value"`
		Score float64 `json:"score"`
	}
	rows := make([]jsonEntry, len(r.scored))
	for i, s := range r.scored {
		rows[i] = jsonEntry{Key: s.Entry.Key, Value: s.Entry.Value, Score: s.Score}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
