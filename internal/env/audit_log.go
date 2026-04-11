package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// AuditLog holds a collection of audit events and provides reporting.
type AuditLog struct {
	Events []AuditEvent
}

// NewAuditLog creates an AuditLog from a slice of events.
func NewAuditLog(events []AuditEvent) *AuditLog {
	return &AuditLog{Events: events}
}

// HasEvents returns true if there are any recorded events.
func (l *AuditLog) HasEvents() bool {
	return len(l.Events) > 0
}

// FilterByAction returns only events matching the given action.
func (l *AuditLog) FilterByAction(action AuditAction) []AuditEvent {
	var result []AuditEvent
	for _, e := range l.Events {
		if e.Action == action {
			result = append(result, e)
		}
	}
	return result
}

// WriteText writes a human-readable audit log to w.
func (l *AuditLog) WriteText(w io.Writer) error {
	if !l.HasEvents() {
		_, err := fmt.Fprintln(w, "No changes detected.")
		return err
	}
	var sb strings.Builder
	for _, e := range l.Events {
		sb.WriteString(e.String())
		sb.WriteByte('\n')
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// WriteJSON writes the audit log as a JSON array to w.
func (l *AuditLog) WriteJSON(w io.Writer) error {
	type jsonEvent struct {
		Timestamp string `json:"timestamp"`
		Action    string `json:"action"`
		Key       string `json:"key"`
		OldValue  string `json:"old_value,omitempty"`
		NewValue  string `json:"new_value,omitempty"`
	}
	var out []jsonEvent
	for _, e := range l.Events {
		out = append(out, jsonEvent{
			Timestamp: e.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			Action:    string(e.Action),
			Key:       e.Key,
			OldValue:  e.OldValue,
			NewValue:  e.NewValue,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
