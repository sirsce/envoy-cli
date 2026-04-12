package env

import (
	"testing"
)

func newTestLabeler() *Labeler {
	return NewLabeler()
}

func TestLabeler_SetAndGet(t *testing.T) {
	l := newTestLabeler()
	if err := l.Set("DB_HOST", "env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ls := l.Get("DB_HOST")
	if ls["env"] != "production" {
		t.Errorf("expected 'production', got %q", ls["env"])
	}
}

func TestLabeler_GetMissing(t *testing.T) {
	l := newTestLabeler()
	ls := l.Get("MISSING_KEY")
	if len(ls) != 0 {
		t.Errorf("expected empty LabelSet, got %v", ls)
	}
}

func TestLabeler_SetEmptyEntryKey(t *testing.T) {
	l := newTestLabeler()
	if err := l.Set("", "env", "production"); err == nil {
		t.Error("expected error for empty entry key")
	}
}

func TestLabeler_SetEmptyLabelKey(t *testing.T) {
	l := newTestLabeler()
	if err := l.Set("DB_HOST", "", "production"); err == nil {
		t.Error("expected error for empty label key")
	}
}

func TestLabeler_Remove(t *testing.T) {
	l := newTestLabeler()
	_ = l.Set("API_KEY", "tier", "secret")
	l.Remove("API_KEY", "tier")
	ls := l.Get("API_KEY")
	if len(ls) != 0 {
		t.Errorf("expected label to be removed, got %v", ls)
	}
}

func TestLabeler_Remove_CleansUpEntry(t *testing.T) {
	l := newTestLabeler()
	_ = l.Set("API_KEY", "tier", "secret")
	l.Remove("API_KEY", "tier")
	if len(l.Keys()) != 0 {
		t.Errorf("expected no labeled keys after removal, got %v", l.Keys())
	}
}

func TestLabeler_FilterByLabel(t *testing.T) {
	l := newTestLabeler()
	_ = l.Set("DB_HOST", "env", "production")
	_ = l.Set("DB_PORT", "env", "production")
	_ = l.Set("DEBUG", "env", "development")

	result := l.FilterByLabel("env", "production")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "DB_HOST" || result[1] != "DB_PORT" {
		t.Errorf("unexpected results: %v", result)
	}
}

func TestLabeler_Keys(t *testing.T) {
	l := newTestLabeler()
	_ = l.Set("Z_KEY", "x", "1")
	_ = l.Set("A_KEY", "x", "1")
	keys := l.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" {
		t.Errorf("expected sorted keys, first is %q", keys[0])
	}
}

func TestLabeler_GetReturnsCopy(t *testing.T) {
	l := newTestLabeler()
	_ = l.Set("DB_HOST", "env", "production")
	ls := l.Get("DB_HOST")
	ls["env"] = "mutated"
	original := l.Get("DB_HOST")
	if original["env"] != "production" {
		t.Error("Get should return a copy, not a reference")
	}
}
