package env

import (
	"fmt"
	"sort"
	"time"
)

// AuditAction represents the type of change made to an env entry.
type AuditAction string

const (
	AuditActionAdd    AuditAction = "add"
	AuditActionRemove AuditAction = "remove"
	AuditActionChange AuditAction = "change"
)

// AuditEvent records a single change to an env entry.
type AuditEvent struct {
	Timestamp time.Time
	Action    AuditAction
	Key       string
	OldValue  string
	NewValue  string
}

func (e AuditEvent) String() string {
	switch e.Action {
	case AuditActionAdd:
		return fmt.Sprintf("[%s] ADD %s=%q", e.Timestamp.Format(time.RFC3339), e.Key, e.NewValue)
	case AuditActionRemove:
		return fmt.Sprintf("[%s] REMOVE %s (was %q)", e.Timestamp.Format(time.RFC3339), e.Key, e.OldValue)
	case AuditActionChange:
		return fmt.Sprintf("[%s] CHANGE %s: %q -> %q", e.Timestamp.Format(time.RFC3339), e.Key, e.OldValue, e.NewValue)
	}
	return ""
}

// Auditor tracks changes between two sets of env entries.
type Auditor struct {
	now func() time.Time
}

// NewAuditor creates a new Auditor.
func NewAuditor() *Auditor {
	return &Auditor{now: time.Now}
}

// Audit compares base and updated entry slices and returns audit events.
// Events are returned in a deterministic order, sorted by key then action.
func (a *Auditor) Audit(base, updated []Entry) []AuditEvent {
	baseMap := toEntryMap(base)
	updatedMap := toEntryMap(updated)

	var events []AuditEvent
	ts := a.now()

	for key, newVal := range updatedMap {
		if oldVal, exists := baseMap[key]; !exists {
			events = append(events, AuditEvent{Timestamp: ts, Action: AuditActionAdd, Key: key, NewValue: newVal})
		} else if oldVal != newVal {
			events = append(events, AuditEvent{Timestamp: ts, Action: AuditActionChange, Key: key, OldValue: oldVal, NewValue: newVal})
		}
	}

	for key, oldVal := range baseMap {
		if _, exists := updatedMap[key]; !exists {
			events = append(events, AuditEvent{Timestamp: ts, Action: AuditActionRemove, Key: key, OldValue: oldVal})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].Key != events[j].Key {
			return events[i].Key < events[j].Key
		}
		return events[i].Action < events[j].Action
	})

	return events
}

func toEntryMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
