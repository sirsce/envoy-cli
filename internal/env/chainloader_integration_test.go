package env_test

import (
	"strings"
	"testing"

	"github.com/yourusername/envoy-cli/internal/env"
)

func TestChainLoader_Integration_ThreeLayers(t *testing.T) {
	base := strings.NewReader("APP_ENV=development\nDB_HOST=localhost\nDEBUG=false\n")
	shared := strings.NewReader("DB_HOST=shared-db\nLOG_LEVEL=info\n")
	local := strings.NewReader("DEBUG=true\nSECRET=abc123\n")

	cl := env.NewChainLoader()
	cl.Add(".env", base)
	cl.Add(".env.shared", shared)
	cl.Add(".env.local", local)

	result, err := cl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Loaded) != 3 {
		t.Errorf("expected 3 loaded sources, got %d: %v", len(result.Loaded), result.Loaded)
	}

	m := make(map[string]string)
	for _, e := range result.Entries {
		m[e.Key] = e.Value
	}

	cases := map[string]string{
		"APP_ENV":   "development",
		"DB_HOST":   "shared-db",
		"DEBUG":     "true",
		"LOG_LEVEL": "info",
		"SECRET":    "abc123",
	}
	for k, want := range cases {
		if got := m[k]; got != want {
			t.Errorf("key %s: want %q, got %q", k, want, got)
		}
	}
}

func TestChainLoader_Integration_WithExporter(t *testing.T) {
	cl := env.NewChainLoader()
	cl.Add("base", strings.NewReader("FOO=1\nBAR=2\n"))
	cl.Add("override", strings.NewReader("FOO=10\nBAZ=3\n"))

	result, err := cl.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exporter := env.NewExporter(result.Entries)
	var buf strings.Builder
	if err := exporter.Write(&buf, env.FormatDotenv); err != nil {
		t.Fatalf("exporter write: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "FOO=10") {
		t.Errorf("expected FOO=10 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAR=2") {
		t.Errorf("expected BAR=2 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=3") {
		t.Errorf("expected BAZ=3 in output, got:\n%s", out)
	}
}
