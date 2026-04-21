package env

import (
	"testing"
)

func makeFlattenerEntries() []Entry {
	return []Entry{
		{Key: "db.host", Value: "localhost"},
		{Key: "db.port", Value: "5432"},
		{Key: "app.name", Value: "envoy"},
		{Key: "PLAIN_KEY", Value: "value"},
	}
}

func TestFlattener_DefaultSeparator(t *testing.T) {
	f := NewFlattener(FlattenOptions{})
	out := f.Flatten(makeFlattenerEntries())

	expected := map[string]string{
		"db_host":  "localhost",
		"db_port":  "5432",
		"app_name": "envoy",
		"PLAIN_KEY": "value",
	}
	for _, e := range out {
		v, ok := expected[e.Key]
		if !ok {
			t.Errorf("unexpected key %q", e.Key)
			continue
		}
		if e.Value != v {
			t.Errorf("key %q: want value %q, got %q", e.Key, v, e.Value)
		}
	}
}

func TestFlattener_Uppercase(t *testing.T) {
	f := NewFlattener(FlattenOptions{Uppercase: true})
	out := f.Flatten([]Entry{{Key: "db.host", Value: "localhost"}})
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %q", out[0].Key)
	}
}

func TestFlattener_WithPrefix(t *testing.T) {
	f := NewFlattener(FlattenOptions{Prefix: "APP"})
	out := f.Flatten([]Entry{{Key: "db.host", Value: "localhost"}})
	if out[0].Key != "APP_db_host" {
		t.Errorf("expected APP_db_host, got %q", out[0].Key)
	}
}

func TestFlattener_CustomSeparator(t *testing.T) {
	f := NewFlattener(FlattenOptions{Separator: ":"})
	out := f.Flatten([]Entry{{Key: "db:host", Value: "localhost"}})
	if out[0].Key != "db_host" {
		t.Errorf("expected db_host, got %q", out[0].Key)
	}
}

func TestFlattener_FlattenedCount(t *testing.T) {
	f := NewFlattener(FlattenOptions{})
	count := f.FlattenedCount(makeFlattenerEntries())
	if count != 3 {
		t.Errorf("expected 3 flattened keys, got %d", count)
	}
}

func TestFlattener_DoesNotMutateOriginal(t *testing.T) {
	original := makeFlattenerEntries()
	f := NewFlattener(FlattenOptions{})
	_ = f.Flatten(original)
	if original[0].Key != "db.host" {
		t.Errorf("original entry was mutated")
	}
}
