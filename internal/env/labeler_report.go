package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// LabelReportEntry represents a single entry in a label report.
type LabelReportEntry struct {
	EntryKey string            `json:"entry_key"`
	Labels   map[string]string `json:"labels"`
}

// LabelReport holds the full label report output.
type LabelReport struct {
	entries []LabelReportEntry
}

// NewLabelReport builds a LabelReport from a Labeler instance.
func NewLabelReport(l *Labeler) *LabelReport {
	keys := l.Keys()
	entries := make([]LabelReportEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, LabelReportEntry{
			EntryKey: k,
			Labels:   l.Get(k),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].EntryKey < entries[j].EntryKey
	})
	return &LabelReport{entries: entries}
}

// WriteText writes a human-readable label report to w.
func (r *LabelReport) WriteText(w io.Writer) error {
	if len(r.entries) == 0 {
		_, err := fmt.Fprintln(w, "No labels defined.")
		return err
	}
	for _, e := range r.entries {
		_, err := fmt.Fprintf(w, "%s:\n", e.EntryKey)
		if err != nil {
			return err
		}
		labelKeys := make([]string, 0, len(e.Labels))
		for k := range e.Labels {
			labelKeys = append(labelKeys, k)
		}
		sort.Strings(labelKeys)
		for _, lk := range labelKeys {
			_, err := fmt.Fprintf(w, "  %s=%s\n", lk, e.Labels[lk])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteJSON writes the label report as JSON to w.
func (r *LabelReport) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if len(r.entries) == 0 {
		return enc.Encode([]LabelReportEntry{})
	}
	return enc.Encode(r.entries)
}
