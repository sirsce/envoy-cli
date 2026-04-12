package env

import (
	"sort"
	"strings"
)

// Group holds a named collection of entries sharing a common prefix.
type Group struct {
	Name    string
	Entries []Entry
}

// Grouper partitions env entries into named groups based on key prefixes.
type Grouper struct {
	separator string
	defaultGroup string
}

// GrouperOption configures a Grouper.
type GrouperOption func(*Grouper)

// WithGroupSeparator sets the separator used to detect prefixes (default "_").
func WithGroupSeparator(sep string) GrouperOption {
	return func(g *Grouper) {
		g.separator = sep
	}
}

// WithDefaultGroupName sets the name used for entries with no detected prefix.
func WithDefaultGroupName(name string) GrouperOption {
	return func(g *Grouper) {
		g.defaultGroup = name
	}
}

// NewGrouper creates a Grouper with the given options.
func NewGrouper(opts ...GrouperOption) *Grouper {
	g := &Grouper{
		separator:    "_",
		defaultGroup: "default",
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

// Group partitions entries by the first segment of their key before the separator.
// Entries without a separator are placed into the default group.
func (g *Grouper) Group(entries []Entry) []Group {
	index := make(map[string][]Entry)
	order := []string{}

	for _, e := range entries {
		name := g.defaultGroup
		if idx := strings.Index(e.Key, g.separator); idx > 0 {
			name = e.Key[:idx]
		}
		if _, exists := index[name]; !exists {
			order = append(order, name)
		}
		index[name] = append(index[name], e)
	}

	sort.Strings(order)

	groups := make([]Group, 0, len(order))
	for _, name := range order {
		groups = append(groups, Group{Name: name, Entries: index[name]})
	}
	return groups
}

// GroupNames returns the sorted list of group names for the given entries.
func (g *Grouper) GroupNames(entries []Entry) []string {
	groups := g.Group(entries)
	names := make([]string, len(groups))
	for i, grp := range groups {
		names[i] = grp.Name
	}
	return names
}
