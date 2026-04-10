package env

import (
	"strings"
	"testing"
)

func TestExporter_FormatDotenv(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	exporter := NewExporter(FormatDotenv, false)
	var buf strings.Builder
	if err := exporter.Export(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO=bar\n") {
		t.Errorf("expected FOO=bar in output, got: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux in output, got: %s", out)
	}
}

func TestExporter_FormatExport(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "secret"},
	}
	exporter := NewExporter(FormatExport, false)
	var buf strings.Builder
	if err := exporter.Export(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export API_KEY=secret\n") {
		t.Errorf("expected 'export API_KEY=secret' in output, got: %s", out)
	}
}

func TestExporter_FormatJSON(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
	exporter := NewExporter(FormatJSON, true)
	var buf strings.Builder
	if err := exporter.Export(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"HOST": "localhost"`) {
		t.Errorf("expected HOST in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"PORT": "8080"`) {
		t.Errorf("expected PORT in JSON output, got: %s", out)
	}
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected valid JSON object braces, got: %s", out)
	}
}

func TestExporter_SortedOutput(t *testing.T) {
	entries := []Entry{
		{Key: "ZEBRA", Value: "1"},
		{Key: "ALPHA", Value: "2"},
		{Key: "MANGO", Value: "3"},
	}
	exporter := NewExporter(FormatDotenv, true)
	var buf strings.Builder
	if err := exporter.Export(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA") {
		t.Errorf("expected ALPHA first, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "MANGO") {
		t.Errorf("expected MANGO second, got: %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected ZEBRA third, got: %s", lines[2])
	}
}

func TestExporter_JSONEscapesQuotes(t *testing.T) {
	entries := []Entry{
		{Key: "MSG", Value: `say "hello"`},
	}
	exporter := NewExporter(FormatJSON, false)
	var buf strings.Builder
	if err := exporter.Export(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"MSG": "say \"hello\""`) {
		t.Errorf("expected escaped quotes in JSON output, got: %s", out)
	}
}
