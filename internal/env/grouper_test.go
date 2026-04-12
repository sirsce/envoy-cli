package env

import (
	"testing"
)

func makeGrouperEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DEBUG", Value: "false"},
	}
}

func TestGrouper_BasicGrouping(t *testing.T) {
	g := NewGrouper()
	groups := g.Group(makeGrouperEntries())

	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}

	names := g.GroupNames(makeGrouperEntries())
	if names[0] != "APP" || names[1] != "DB" || names[2] != "default" {
		t.Errorf("unexpected group names: %v", names)
	}
}

func TestGrouper_DefaultGroup(t *testing.T) {
	g := NewGrouper()
	groups := g.Group(makeGrouperEntries())

	var def *Group
	for i := range groups {
		if groups[i].Name == "default" {
			def = &groups[i]
			break
		}
	}
	if def == nil {
		t.Fatal("expected default group")
	}
	if len(def.Entries) != 1 || def.Entries[0].Key != "DEBUG" {
		t.Errorf("unexpected default group entries: %v", def.Entries)
	}
}

func TestGrouper_CustomSeparator(t *testing.T) {
	entries := []Entry{
		{Key: "DB.HOST", Value: "localhost"},
		{Key: "DB.PORT", Value: "5432"},
		{Key: "PLAIN", Value: "val"},
	}
	g := NewGrouper(WithGroupSeparator("."))
	groups := g.Group(entries)

	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}

func TestGrouper_CustomDefaultName(t *testing.T) {
	entries := []Entry{
		{Key: "STANDALONE", Value: "yes"},
	}
	g := NewGrouper(WithDefaultGroupName("ungrouped"))
	groups := g.Group(entries)

	if len(groups) != 1 || groups[0].Name != "ungrouped" {
		t.Errorf("expected ungrouped, got %v", groups)
	}
}

func TestGrouper_Empty(t *testing.T) {
	g := NewGrouper()
	groups := g.Group([]Entry{})
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestGrouper_SortedOutput(t *testing.T) {
	entries := []Entry{
		{Key: "Z_ONE", Value: "1"},
		{Key: "A_TWO", Value: "2"},
		{Key: "M_THREE", Value: "3"},
	}
	g := NewGrouper()
	names := g.GroupNames(entries)
	if names[0] != "A" || names[1] != "M" || names[2] != "Z" {
		t.Errorf("expected sorted group names, got %v", names)
	}
}
