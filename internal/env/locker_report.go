package env

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// LockReport renders a human-readable or JSON report of active locks.
type LockReport struct {
	entries []LockEntry
}

// NewLockReport creates a LockReport from a slice of LockEntry values.
func NewLockReport(entries []LockEntry) *LockReport {
	sorted := make([]LockEntry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	return &LockReport{entries: sorted}
}

// WriteText writes a plain-text summary of locks to w.
func (r *LockReport) WriteText(w io.Writer) error {
	if len(r.entries) == 0 {
		_, err := fmt.Fprintln(w, "No active locks.")
		return err
	}
	for _, e := range r.entries {
		expiry := "never"
		if e.ExpiresAt != nil {
			expiry = e.ExpiresAt.Format("2006-01-02T15:04:05Z")
		}
		_, err := fmt.Fprintf(w, "[LOCKED] %s  owner=%s  locked_at=%s  expires=%s\n",
			e.Key, e.LockedBy, e.LockedAt.Format("2006-01-02T15:04:05Z"), expiry)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON writes a JSON array of lock entries to w.
func (r *LockReport) WriteJSON(w io.Writer) error {
	type jsonEntry struct {
		Key       string  `json:"key"`
		LockedBy  string  `json:"locked_by"`
		LockedAt  string  `json:"locked_at"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	out := make([]jsonEntry, len(r.entries))
	for i, e := range r.entries {
		je := jsonEntry{
			Key:      e.Key,
			LockedBy: e.LockedBy,
			LockedAt: e.LockedAt.Format("2006-01-02T15:04:05Z"),
		}
		if e.ExpiresAt != nil {
			s := e.ExpiresAt.Format("2006-01-02T15:04:05Z")
			je.ExpiresAt = &s
		}
		out[i] = je
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
