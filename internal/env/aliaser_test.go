package env

import (
	"testing"
)

func makeAliasEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestAliaser_CreatesAlias(t *testing.T) {
	a, err := NewAliaser([]AliasRule{{Alias: "DATABASE_HOST", Source: "DB_HOST"}}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, results := a.Apply(makeAliasEntries())
	if AppliedAliasCount(results) != 1 {
		t.Fatalf("expected 1 applied, got %d", AppliedAliasCount(results))
	}
	found := false
	for _, e := range out {
		if e.Key == "DATABASE_HOST" && e.Value == "localhost" {
			found = true
		}
	}
	if !found {
		t.Error("alias entry DATABASE_HOST not found in output")
	}
}

func TestAliaser_SourceNotFound(t *testing.T) {
	a, _ := NewAliaser([]AliasRule{{Alias: "MISSING_ALIAS", Source: "NONEXISTENT"}}, false)
	_, results := a.Apply(makeAliasEntries())
	if len(results) != 1 || results[0].Applied {
		t.Error("expected one unapplied result for missing source")
	}
	if results[0].Reason != "source key not found" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestAliaser_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := append(makeAliasEntries(), Entry{Key: "DATABASE_HOST", Value: "old"})
	a, _ := NewAliaser([]AliasRule{{Alias: "DATABASE_HOST", Source: "DB_HOST"}}, false)
	out, results := a.Apply(entries)
	if results[0].Applied {
		t.Error("should not overwrite existing alias key when overwrite=false")
	}
	for _, e := range out {
		if e.Key == "DATABASE_HOST" && e.Value != "old" {
			t.Error("value should remain 'old' without overwrite")
		}
	}
}

func TestAliaser_OverwritesExisting(t *testing.T) {
	entries := append(makeAliasEntries(), Entry{Key: "DATABASE_HOST", Value: "old"})
	a, _ := NewAliaser([]AliasRule{{Alias: "DATABASE_HOST", Source: "DB_HOST"}}, true)
	out, results := a.Apply(entries)
	if !results[0].Applied {
		t.Error("expected alias to be applied with overwrite=true")
	}
	for _, e := range out {
		if e.Key == "DATABASE_HOST" && e.Value != "localhost" {
			t.Errorf("expected overwritten value 'localhost', got %q", e.Value)
		}
	}
}

func TestAliaser_InvalidRule_EmptyAlias(t *testing.T) {
	_, err := NewAliaser([]AliasRule{{Alias: "", Source: "DB_HOST"}}, false)
	if err == nil {
		t.Error("expected error for empty alias key")
	}
}

func TestAliaser_InvalidRule_SameAliasAndSource(t *testing.T) {
	_, err := NewAliaser([]AliasRule{{Alias: "DB_HOST", Source: "DB_HOST"}}, false)
	if err == nil {
		t.Error("expected error when alias equals source")
	}
}

func TestAliaser_MultipleRules(t *testing.T) {
	rules := []AliasRule{
		{Alias: "HOST", Source: "DB_HOST"},
		{Alias: "PORT", Source: "DB_PORT"},
	}
	a, _ := NewAliaser(rules, false)
	out, results := a.Apply(makeAliasEntries())
	if AppliedAliasCount(results) != 2 {
		t.Fatalf("expected 2 applied, got %d", AppliedAliasCount(results))
	}
	if len(out) != len(makeAliasEntries())+2 {
		t.Errorf("expected %d entries, got %d", len(makeAliasEntries())+2, len(out))
	}
}
