package env

import (
	"testing"
)

func makeClonerEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestCloner_ClonesKey(t *testing.T) {
	c := NewCloner(makeClonerEntries(), false)
	if err := c.Clone("DB_HOST", "DB_HOST_BACKUP"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := c.Entries()
	for _, e := range entries {
		if e.Key == "DB_HOST_BACKUP" && e.Value == "localhost" {
			return
		}
	}
	t.Fatal("expected DB_HOST_BACKUP with value 'localhost'")
}

func TestCloner_SourceNotFound(t *testing.T) {
	c := NewCloner(makeClonerEntries(), false)
	err := c.Clone("MISSING_KEY", "DEST")
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestCloner_SkipsExistingWithoutOverwrite(t *testing.T) {
	c := NewCloner(makeClonerEntries(), false)
	if err := c.Clone("DB_HOST", "DB_PORT"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := c.Results()
	if len(results) != 1 || !results[0].Skipped {
		t.Fatal("expected clone to be skipped")
	}
	for _, e := range c.Entries() {
		if e.Key == "DB_PORT" && e.Value != "5432" {
			t.Fatal("expected DB_PORT to remain unchanged")
		}
	}
}

func TestCloner_OverwritesExisting(t *testing.T) {
	c := NewCloner(makeClonerEntries(), true)
	if err := c.Clone("DB_HOST", "DB_PORT"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range c.Entries() {
		if e.Key == "DB_PORT" && e.Value == "localhost" {
			return
		}
	}
	t.Fatal("expected DB_PORT to be overwritten with 'localhost'")
}

func TestCloner_EmptyKeyReturnsError(t *testing.T) {
	c := NewCloner(makeClonerEntries(), false)
	if err := c.Clone("", "DEST"); err == nil {
		t.Fatal("expected error for empty source key")
	}
	if err := c.Clone("DB_HOST", ""); err == nil {
		t.Fatal("expected error for empty dest key")
	}
}

func TestCloner_DoesNotMutateOriginal(t *testing.T) {
	orig := makeClonerEntries()
	c := NewCloner(orig, false)
	_ = c.Clone("DB_HOST", "DB_HOST_COPY")
	if len(orig) != 3 {
		t.Fatal("original slice should not be mutated")
	}
}

func TestCloner_AppliedCount(t *testing.T) {
	c := NewCloner(makeClonerEntries(), false)
	_ = c.Clone("DB_HOST", "DB_HOST_COPY")
	_ = c.Clone("DB_PORT", "DB_HOST_COPY") // skipped, already exists
	if c.AppliedCount() != 1 {
		t.Fatalf("expected AppliedCount 1, got %d", c.AppliedCount())
	}
}
