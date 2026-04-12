package env

import (
	"testing"
)

func newTestTagger() *Tagger {
	return NewTagger()
}

func TestTagger_TagAndHasTag(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("DB_HOST", "sensitive", "")

	if !tg.HasTag("DB_HOST", "sensitive") {
		t.Error("expected DB_HOST to have tag 'sensitive'")
	}
	if tg.HasTag("DB_HOST", "optional") {
		t.Error("expected DB_HOST NOT to have tag 'optional'")
	}
}

func TestTagger_GetTags(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("API_KEY", "sensitive", "")
	tg.Tag("API_KEY", "env", "production")

	tags := tg.GetTags("API_KEY")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Name != "sensitive" {
		t.Errorf("expected first tag 'sensitive', got %q", tags[0].Name)
	}
	if tags[1].Name != "env" || tags[1].Value != "production" {
		t.Errorf("unexpected second tag: %+v", tags[1])
	}
}

func TestTagger_FilterByTag(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("DB_PASS", "sensitive", "")
	tg.Tag("API_KEY", "sensitive", "")
	tg.Tag("APP_PORT", "optional", "")

	entries := []Entry{
		{Key: "DB_PASS", Value: "secret"},
		{Key: "API_KEY", Value: "abc123"},
		{Key: "APP_PORT", Value: "8080"},
	}

	result := tg.FilterByTag(entries, "sensitive")
	if len(result) != 2 {
		t.Fatalf("expected 2 sensitive entries, got %d", len(result))
	}
}

func TestTagger_RemoveTag(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("DB_HOST", "sensitive", "")
	tg.Tag("DB_HOST", "required", "")

	tg.RemoveTag("DB_HOST", "sensitive")

	if tg.HasTag("DB_HOST", "sensitive") {
		t.Error("expected 'sensitive' tag to be removed")
	}
	if !tg.HasTag("DB_HOST", "required") {
		t.Error("expected 'required' tag to still exist")
	}
}

func TestTagger_Keys(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("Z_KEY", "foo", "")
	tg.Tag("A_KEY", "bar", "")
	tg.Tag("M_KEY", "baz", "")

	keys := tg.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("keys not sorted: %v", keys)
	}
}

func TestTagger_Summary(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("DB_HOST", "sensitive", "")
	tg.Tag("DB_HOST", "env", "prod")

	summary := tg.Summary("DB_HOST")
	if summary == "" {
		t.Error("expected non-empty summary")
	}

	empty := tg.Summary("UNKNOWN")
	if empty != "UNKNOWN: (no tags)" {
		t.Errorf("unexpected empty summary: %q", empty)
	}
}

func TestTagger_IgnoresEmptyKeyOrName(t *testing.T) {
	tg := newTestTagger()
	tg.Tag("", "sensitive", "")
	tg.Tag("DB_HOST", "", "")

	if len(tg.Keys()) != 0 {
		t.Error("expected no keys to be registered for empty key/name")
	}
}
