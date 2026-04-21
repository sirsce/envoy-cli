package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestCaster_Integration_WithParser(t *testing.T) {
	raw := `PORT= 9000 
DEBUG=true
RATE= 2.71 
NAME=envoy
`
	p := NewParser()
	entries, err := p.Parse(bytes.NewBufferString(raw))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	c := NewCaster([]CastRule{
		{Key: "PORT", CastTo: CastInt},
		{Key: "DEBUG", CastTo: CastBool},
		{Key: "RATE", CastTo: CastFloat},
	})
	out, results := c.Apply(entries)

	if CastedCount(results) != 3 {
		t.Errorf("expected 3 casted, got %d", CastedCount(results))
	}

	m := make(map[string]string)
	for _, e := range out {
		m[e.Key] = e.Value
	}
	if m["PORT"] != "9000" {
		t.Errorf("PORT: got %q", m["PORT"])
	}
	if m["DEBUG"] != "true" {
		t.Errorf("DEBUG: got %q", m["DEBUG"])
	}
	if m["RATE"] != "2.71" {
		t.Errorf("RATE: got %q", m["RATE"])
	}
	if m["NAME"] != "envoy" {
		t.Errorf("NAME should be unchanged, got %q", m["NAME"])
	}
}

func TestCaster_Integration_ReportAndExport(t *testing.T) {
	entries := []Entry{
		{Key: "TIMEOUT", Value: " 30 "},
		{Key: "VERBOSE", Value: "0"},
		{Key: "BROKEN", Value: "xyz"},
	}

	c := NewCaster([]CastRule{
		{Key: "TIMEOUT", CastTo: CastInt},
		{Key: "VERBOSE", CastTo: CastBool},
		{Key: "BROKEN", CastTo: CastFloat},
	})
	out, results := c.Apply(entries)

	rep := NewCastReport(results)
	var buf bytes.Buffer
	if err := rep.WriteText(&buf); err != nil {
		t.Fatalf("report error: %v", err)
	}
	if !strings.Contains(buf.String(), "TIMEOUT") {
		t.Error("expected TIMEOUT in report")
	}

	ex := NewExporter()
	var out2 bytes.Buffer
	if err := ex.Export(out, &out2, FormatDotenv); err != nil {
		t.Fatalf("export error: %v", err)
	}
	if !strings.Contains(out2.String(), "TIMEOUT=30") {
		t.Errorf("expected TIMEOUT=30 in export, got:\n%s", out2.String())
	}
}
