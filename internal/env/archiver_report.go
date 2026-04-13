package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// ArchiveReport renders archive listings as text or JSON.
type ArchiveReport struct {
	entries []ArchiveEntry
}

// NewArchiveReport creates a report from the given archive entries.
func NewArchiveReport(entries []ArchiveEntry) *ArchiveReport {
	return &ArchiveReport{entries: entries}
}

// WriteText writes a human-readable archive listing to w.
func (r *ArchiveReport) WriteText(w io.Writer) error {
	if len(r.entries) == 0 {
		_, err := fmt.Fprintln(w, "No archives found.")
		return err
	}
	for _, ae := range r.entries {
		_, err := fmt.Fprintf(w, "v%-4d  %-20s  %d keys  %s\n",
			ae.Version,
			ae.Label,
			len(ae.Entries),
			ae.ArchivedAt.Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type archiveReportJSON struct {
	Version    int    `json:"version"`
	Label      string `json:"label"`
	KeyCount   int    `json:"key_count"`
	ArchivedAt string `json:"archived_at"`
}

// WriteJSON writes a JSON array of archive metadata to w.
func (r *ArchiveReport) WriteJSON(w io.Writer) error {
	out := make([]archiveReportJSON, len(r.entries))
	for i, ae := range r.entries {
		out[i] = archiveReportJSON{
			Version:    ae.Version,
			Label:      ae.Label,
			KeyCount:   len(ae.Entries),
			ArchivedAt: ae.ArchivedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
