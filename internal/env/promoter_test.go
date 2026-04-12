package env

import (
	"testing"
)

func makePromoterScopes() map[string][]Entry {
	return map[string][]Entry{
		"dev": {
			{Key: "APP_URL", Value: "http://localhost"},
			{Key: "DB_HOST", Value: "localhost"},
		},
		"staging": {
			{Key: "APP_URL", Value: "http://staging.example.com"},
		},
	}
}

func TestPromoter_NewScopeCreated(t *testing.T) {
	p := NewPromoter(makePromoterScopes(), false)
	rule := PromoteRule{From: "dev", To: "prod"}
	results, err := p.Promote(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	entries, ok := p.Scope("prod")
	if !ok {
		t.Fatal("expected prod scope to exist")
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries in prod, got %d", len(entries))
	}
}

func TestPromoter_SkipsExistingWithoutOverwrite(t *testing.T) {
	p := NewPromoter(makePromoterScopes(), false)
	rule := PromoteRule{From: "dev", To: "staging"}
	results, err := p.Promote(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var skipped int
	for _, r := range results {
		if r.Skipped {
			skipped++
		}
	}
	if skipped != 1 {
		t.Fatalf("expected 1 skipped, got %d", skipped)
	}
}

func TestPromoter_OverwritesExisting(t *testing.T) {
	p := NewPromoter(makePromoterScopes(), true)
	rule := PromoteRule{From: "dev", To: "staging"}
	results, err := p.Promote(rule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var overwrote int
	for _, r := range results {
		if r.Overwrote {
			overwrote++
		}
	}
	if overwrote != 1 {
		t.Fatalf("expected 1 overwrite, got %d", overwrote)
	}
	entries, _ := p.Scope("staging")
	for _, e := range entries {
		if e.Key == "APP_URL" && e.Value != "http://localhost" {
			t.Errorf("expected overwritten value, got %q", e.Value)
		}
	}
}

func TestPromoter_SourceNotFound(t *testing.T) {
	p := NewPromoter(makePromoterScopes(), false)
	rule := PromoteRule{From: "nonexistent", To: "staging"}
	_, err := p.Promote(rule)
	if err == nil {
		t.Fatal("expected error for missing source scope")
	}
}

func TestPromoter_DoesNotMutateOriginalSlice(t *testing.T) {
	scopes := makePromoterScopes()
	origLen := len(scopes["dev"])
	p := NewPromoter(scopes, false)
	p.Promote(PromoteRule{From: "dev", To: "staging"})
	if len(scopes["dev"]) != origLen {
		t.Error("source scope was mutated")
	}
}
