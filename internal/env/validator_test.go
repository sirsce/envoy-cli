package env

import (
	"testing"
)

func TestValidator_ValidEntries(t *testing.T) {
	v := NewValidator()
	entries := []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "PORT", Value: "8080"},
		{Key: "_PRIVATE", Value: "secret"},
	}
	result := v.Validate(entries)
	if !result.IsValid() {
		t.Errorf("expected valid, got errors: %s", result.Error())
	}
}

func TestValidator_InvalidKeyName(t *testing.T) {
	v := NewValidator()
	entries := []Entry{
		{Key: "123BAD", Value: "value"},
		{Key: "GOOD_KEY", Value: "ok"},
		{Key: "has-hyphen", Value: "nope"},
	}
	result := v.Validate(entries)
	if result.IsValid() {
		t.Fatal("expected validation errors for invalid key names")
	}
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %s", len(result.Errors), result.Error())
	}
}

func TestValidator_DuplicateKey(t *testing.T) {
	v := NewValidator()
	entries := []Entry{
		{Key: "FOO", Value: "first"},
		{Key: "BAR", Value: "bar"},
		{Key: "FOO", Value: "second"},
	}
	result := v.Validate(entries)
	if result.IsValid() {
		t.Fatal("expected duplicate key error")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Key != "FOO" {
		t.Errorf("expected error for key FOO, got %q", result.Errors[0].Key)
	}
}

func TestValidator_RequiredKeys(t *testing.T) {
	v := NewValidator("DATABASE_URL", "SECRET_KEY")
	entries := []Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DATABASE_URL", Value: "postgres://localhost/db"},
	}
	result := v.Validate(entries)
	if result.IsValid() {
		t.Fatal("expected missing required key error")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d: %s", len(result.Errors), result.Error())
	}
	if result.Errors[0].Key != "SECRET_KEY" {
		t.Errorf("expected error for SECRET_KEY, got %q", result.Errors[0].Key)
	}
}

func TestValidator_AllRequiredPresent(t *testing.T) {
	v := NewValidator("A", "B")
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	result := v.Validate(entries)
	if !result.IsValid() {
		t.Errorf("expected valid, got: %s", result.Error())
	}
}
