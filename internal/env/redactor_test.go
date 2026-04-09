package env

import (
	"testing"
)

func TestRedactor_IsSensitive(t *testing.T) {
	r := NewRedactor(RedactFull)

	sensitiveKeys := []string{
		"PASSWORD", "DB_PASSWORD", "API_KEY", "SECRET",
		"AUTH_TOKEN", "PRIVATE_KEY", "aws_secret", "user_token",
	}
	for _, key := range sensitiveKeys {
		if !r.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}

	insensitiveKeys := []string{"PORT", "HOST", "DEBUG", "APP_NAME", "LOG_LEVEL"}
	for _, key := range insensitiveKeys {
		if r.IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestRedactor_ExtraKeys(t *testing.T) {
	r := NewRedactor(RedactFull, "CUSTOM_FIELD", "MY_SECRET_STUFF")

	if !r.IsSensitive("CUSTOM_FIELD") {
		t.Error("expected CUSTOM_FIELD to be sensitive")
	}
	if !r.IsSensitive("custom_field") {
		t.Error("expected custom_field (lowercase) to be sensitive")
	}
	if r.IsSensitive("NORMAL_KEY") {
		t.Error("expected NORMAL_KEY to NOT be sensitive")
	}
}

func TestRedactor_RedactFull(t *testing.T) {
	r := NewRedactor(RedactFull)

	got := r.Redact("PASSWORD", "supersecret")
	if got != "***********" {
		t.Errorf("expected full redaction, got %q", got)
	}

	got = r.Redact("HOST", "localhost")
	if got != "localhost" {
		t.Errorf("expected no redaction for non-sensitive key, got %q", got)
	}
}

func TestRedactor_RedactPartial(t *testing.T) {
	r := NewRedactor(RedactPartial)

	got := r.Redact("API_KEY", "abcdefgh")
	if got != "a******h" {
		t.Errorf("expected partial redaction a******h, got %q", got)
	}

	got = r.Redact("API_KEY", "ab")
	if got != "**" {
		t.Errorf("expected ** for short value, got %q", got)
	}
}

func TestRedactor_RedactEmptyValue(t *testing.T) {
	r := NewRedactor(RedactFull)
	got := r.Redact("PASSWORD", "")
	if got != "" {
		t.Errorf("expected empty string for empty value, got %q", got)
	}
}

func TestRedactor_RedactEntries(t *testing.T) {
	r := NewRedactor(RedactFull)
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Key: "PORT", Value: "5432"},
		{Key: "API_KEY", Value: "mykey"},
	}

	redacted := r.RedactEntries(entries)

	if redacted[0].Value != "localhost" {
		t.Errorf("HOST should not be redacted")
	}
	if redacted[1].Value != "%s" {
		// check it's all stars
		for _, c := range redacted[1].Value {
			if c != '*' {
				t.Errorf("DB_PASSWORD value should be fully redacted, got %q", redacted[1].Value)
				break
			}
		}
	}
	if redacted[2].Value != "5432" {
		t.Errorf("PORT should not be redacted")
	}
	// ensure original slice is not mutated
	if entries[1].Value != "secret123" {
		t.Error("original entries should not be mutated")
	}
}
