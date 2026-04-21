package env

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func makeEncoderEntries() []Entry {
	return []Entry{
		{Key: "DB_PASSWORD", Value: "secret"},
		{Key: "API_KEY", Value: "mykey"},
		{Key: "APP_NAME", Value: "envoy"},
	}
}

func TestEncoder_Base64_AllKeys(t *testing.T) {
	enc := NewEncoder(EncodeBase64)
	out, err := enc.Encode(makeEncoderEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		dec, err := base64.StdEncoding.DecodeString(e.Value)
		if err != nil {
			t.Errorf("key %q value %q is not valid base64: %v", e.Key, e.Value, err)
		}
		_ = dec
	}
	if enc.EncodedCount() != 3 {
		t.Errorf("expected 3 encoded, got %d", enc.EncodedCount())
	}
}

func TestEncoder_Base64_SpecificKey(t *testing.T) {
	enc := NewEncoder(EncodeBase64, "DB_PASSWORD")
	out, err := enc.Encode(makeEncoderEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enc.EncodedCount() != 1 {
		t.Errorf("expected 1 encoded, got %d", enc.EncodedCount())
	}
	if out[1].Value != "mykey" {
		t.Errorf("expected API_KEY unchanged, got %q", out[1].Value)
	}
	expected := base64.StdEncoding.EncodeToString([]byte("secret"))
	if out[0].Value != expected {
		t.Errorf("expected DB_PASSWORD encoded to %q, got %q", expected, out[0].Value)
	}
}

func TestEncoder_Hex(t *testing.T) {
	enc := NewEncoder(EncodeHex, "APP_NAME")
	out, err := enc.Encode(makeEncoderEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "envoy" in hex = 656e766f79
	if out[2].Value != "656e766f79" {
		t.Errorf("expected hex encoded value, got %q", out[2].Value)
	}
}

func TestEncoder_URL(t *testing.T) {
	enc := NewEncoder(EncodeURL, "API_KEY")
	out, err := enc.Encode(makeEncoderEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dec, err := base64.URLEncoding.DecodeString(out[1].Value)
	if err != nil {
		t.Fatalf("invalid url-base64: %v", err)
	}
	if string(dec) != "mykey" {
		t.Errorf("expected decoded value 'mykey', got %q", string(dec))
	}
}

func TestEncoder_DoesNotMutateOriginal(t *testing.T) {
	original := makeEncoderEntries()
	enc := NewEncoder(EncodeBase64)
	_, err := enc.Encode(original)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if original[0].Value != "secret" {
		t.Errorf("original entry was mutated")
	}
}

func TestEncodeReport_TextOutput(t *testing.T) {
	enc := NewEncoder(EncodeBase64, "DB_PASSWORD")
	_, _ = enc.Encode(makeEncoderEntries())
	report := NewEncodeReport(enc.Results())
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD in report output")
	}
}

func TestEncodeReport_JSONOutput(t *testing.T) {
	enc := NewEncoder(EncodeBase64)
	_, _ = enc.Encode(makeEncoderEntries())
	report := NewEncodeReport(enc.Results())
	var buf bytes.Buffer
	if err := report.WriteJSON(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "\"key\"") {
		t.Errorf("expected JSON key field in output")
	}
}

func TestEncodeReport_Empty(t *testing.T) {
	report := NewEncodeReport(nil)
	var buf bytes.Buffer
	if err := report.WriteText(&buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No entries") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}
