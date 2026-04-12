package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// GroupReport renders a summary of grouped env entries.
type GroupReport struct {
	groups []Group
}

// NewGroupReport creates a GroupReport from the provided groups.
func NewGroupReport(groups []Group) *GroupReport {
	return &GroupReport{groups: groups}
}

// WriteText writes a human-readable group report to w.
func (r *GroupReport) WriteText(w io.Writer) error {
	if len(r.groups) == 0 {
		_, err := fmt.Fprintln(w, "No groups found.")
		return err
	}
	for _, g := range r.groups {
		_, err := fmt.Fprintf(w, "[%s] (%d keys)\n", g.Name, len(g.Entries))
		if err != nil {
			return err
		}
		for _, e := range g.Entries {
			_, err = fmt.Fprintf(w, "  %s\n", e.Key)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// jsonGroup is the serialisable form of a Group.
type jsonGroup struct {
	Name    string   `json:"name"`
	Keys    []string `json:"keys"`
}

// WriteJSON writes a JSON-encoded group report to w.
func (r *GroupReport) WriteJSON(w io.Writer) error {
	out := make([]jsonGroup, len(r.groups))
	for i, g := range r.groups {
		keys := make([]string, len(g.Entries))
		for j, e := range g.Entries {
			keys[j] = e.Key
		}
		out[i] = jsonGroup{Name: g.Name, Keys: keys}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
