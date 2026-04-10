package env

import (
	"sort"
	"testing"
)

func TestProfileManager_AddGet(t *testing.T) {
	pm := NewProfileManager()
	entries := []Entry{{Key: "APP_ENV", Value: "development"}}
	pm.Add("dev", entries)

	p, err := pm.Get("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name != "dev" {
		t.Errorf("expected name %q, got %q", "dev", p.Name)
	}
	if len(p.Entries) != 1 || p.Entries[0].Key != "APP_ENV" {
		t.Errorf("unexpected entries: %v", p.Entries)
	}
}

func TestProfileManager_GetMissing(t *testing.T) {
	pm := NewProfileManager()
	_, err := pm.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing profile")
	}
}

func TestProfileManager_List(t *testing.T) {
	pm := NewProfileManager()
	pm.Add("dev", nil)
	pm.Add("prod", nil)

	names := pm.List()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "dev" || names[1] != "prod" {
		t.Errorf("unexpected profiles: %v", names)
	}
}

func TestProfileManager_Merge(t *testing.T) {
	pm := NewProfileManager()
	pm.Add("base", []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DEBUG", Value: "false"},
	})
	pm.Add("prod", []Entry{
		{Key: "DB_HOST", Value: "prod.db.internal"},
		{Key: "APP_ENV", Value: "production"},
	})

	merged, err := pm.Merge("base", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := make(map[string]string)
	for _, e := range merged {
		m[e.Key] = e.Value
	}

	if m["DB_HOST"] != "prod.db.internal" {
		t.Errorf("expected overridden DB_HOST, got %q", m["DB_HOST"])
	}
	if m["DEBUG"] != "false" {
		t.Errorf("expected base DEBUG, got %q", m["DEBUG"])
	}
	if m["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV from prod, got %q", m["APP_ENV"])
	}
}

func TestProfileManager_MergeMissingBase(t *testing.T) {
	pm := NewProfileManager()
	pm.Add("prod", nil)
	_, err := pm.Merge("base", "prod")
	if err == nil {
		t.Fatal("expected error for missing base profile")
	}
}
