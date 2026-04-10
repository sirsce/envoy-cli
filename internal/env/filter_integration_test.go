package env_test

import (
	"strings"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/env"
)

// TestFilter_WithParser ensures Filter works correctly with parsed .env content.
func TestFilter_WithParser(t *testing.T) {
	input := `APP_HOST=localhost
APP_PORT=8080
DB_HOST=postgres
DB_PASSWORD=secret
APP_DEBUG=true
`
	p := env.NewParser()
	entries, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	f := env.NewFilter(env.FilterOptions{Prefix: "APP_"})
	filtered := f.Apply(entries)

	if len(filtered) != 3 {
		t.Fatalf("expected 3 APP_ entries, got %d", len(filtered))
	}
	for _, e := range filtered {
		if !strings.HasPrefix(e.Key, "APP_") {
			t.Errorf("unexpected key %q in filtered result", e.Key)
		}
	}
}

// TestFilter_StripAndExport ensures StripPrefix + Exporter produce clean output.
func TestFilter_StripAndExport(t *testing.T) {
	input := `APP_HOST=localhost
APP_PORT=8080
`
	p := env.NewParser()
	entries, err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	stripped := env.StripPrefix(entries, "APP_")
	if len(stripped) != 2 {
		t.Fatalf("expected 2 entries after strip, got %d", len(stripped))
	}

	keyMap := make(map[string]string, len(stripped))
	for _, e := range stripped {
		keyMap[e.Key] = e.Value
	}
	if keyMap["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", keyMap["HOST"])
	}
	if keyMap["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", keyMap["PORT"])
	}
}
