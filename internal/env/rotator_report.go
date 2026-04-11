package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// RotationReport wraps a slice of RotationRecord and provides rendering helpers.
type RotationReport struct {
	Records []RotationRecord
}

// NewRotationReport creates a RotationReport from the given records.
func NewRotationReport(records []RotationRecord) *RotationReport {
	return &RotationReport{Records: records}
}

// HasRotations reports whether any rotations were recorded.
func (rr *RotationReport) HasRotations() bool {
	return len(rr.Records) > 0
}

// rotatedAtFormat is the timestamp layout used in all report output.
const rotatedAtFormat = "2006-01-02T15:04:05Z"

// WriteText writes a human-readable rotation summary to w.
func (rr *RotationReport) WriteText(w io.Writer) error {
	if !rr.HasRotations() {
		_, err := fmt.Fprintln(w, "No keys rotated.")
		return err
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rotated %d key(s):\n", len(rr.Records)))
	for _, rec := range rr.Records {
		sb.WriteString(fmt.Sprintf("  %-30s  %s\n", rec.Key, rec.RotatedAt.Format(rotatedAtFormat)))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// WriteJSON writes the rotation records as a JSON array to w.
func (rr *RotationReport) WriteJSON(w io.Writer) error {
	type jsonRecord struct {
		Key       string `json:"key"`
		RotatedAt string `json:"rotated_at"`
	}
	out := make([]jsonRecord, len(rr.Records))
	for i, r := range rr.Records {
		out[i] = jsonRecord{
			Key:       r.Key,
			RotatedAt: r.RotatedAt.Format(rotatedAtFormat),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
