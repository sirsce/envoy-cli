package env

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRotator_Integration_WithParser(t *testing.T) {
	raw := []byte("DB_PASSWORD=secret\nAPI_KEY=key123\nAPP_ENV=production\n")
	p := NewParser()
	entries, err := p.ParseBytes(raw)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	counter := 0
	r := newTestRotator()
	r.Register("DB_PASSWORD", func(_ string) (string, error) {
		counter++
		return fmt.Sprintf("rotated-pass-%d", counter), nil
	})
	r.Register("API_KEY", func(old string) (string, error) {
		return old + "-rotated", nil
	})

	out, records, err := r.Rotate(entries)
	if err != nil {
		t.Fatalf("rotate error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	m := make(map[string]string)
	for _, e := range out {
		m[e.Key] = e.Value
	}
	if m["DB_PASSWORD"] != "rotated-pass-1" {
		t.Errorf("unexpected DB_PASSWORD: %q", m["DB_PASSWORD"])
	}
	if m["API_KEY"] != "key123-rotated" {
		t.Errorf("unexpected API_KEY: %q", m["API_KEY"])
	}
	if m["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged")
	}
}

func TestRotator_Integration_ReportAndExport(t *testing.T) {
	r := newTestRotator()
	r.Register("TOKEN", func(_ string) (string, error) { return "new-token", nil })

	entries := []Entry{{Key: "TOKEN", Value: "old-token"}, {Key: "HOST", Value: "localhost"}}
	_, records, err := r.Rotate(entries)
	if err != nil {
		t.Fatalf("rotate error: %v", err)
	}

	report := NewRotationReport(records)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(buf.String(), "TOKEN") {
		t.Errorf("report should mention TOKEN, got: %q", buf.String())
	}
}
