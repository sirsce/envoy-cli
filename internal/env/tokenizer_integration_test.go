package env

import (
	"bytes"
	"strings"
	"testing"
)

const tokenizerSample = `# Database config
DB_HOST=localhost
DB_PORT=5432

# App secrets
APP_SECRET="s3cr3t"
APP_DEBUG=false
`

func TestTokenizer_Integration_FullFile(t *testing.T) {
	tok := NewTokenizer(WithKeepComments(), WithKeepBlanks())
	tokens, err := tok.Tokenize(tokenizerSample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := NewTokenReport(tokens)
	if r.CountByType(TokenKey) != 4 {
		t.Errorf("expected 4 KEY tokens, got %d", r.CountByType(TokenKey))
	}
	if r.CountByType(TokenComment) != 2 {
		t.Errorf("expected 2 COMMENT tokens, got %d", r.CountByType(TokenComment))
	}
	if r.CountByType(TokenBlank) != 1 {
		t.Errorf("expected 1 BLANK token, got %d", r.CountByType(TokenBlank))
	}
}

func TestTokenizer_Integration_ReportAndExport(t *testing.T) {
	tok := NewTokenizer()
	tokens, err := tok.Tokenize(tokenizerSample)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Build entries from KEY tokens for downstream use.
	var entries []Entry
	for _, tk := range tokens {
		if tk.Type == TokenKey {
			entries = append(entries, Entry{Key: tk.Key, Value: tk.Value})
		}
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	// Verify text report contains expected keys.
	r := NewTokenReport(tokens)
	var buf bytes.Buffer
	if err := r.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	for _, key := range []string{"DB_HOST", "DB_PORT", "APP_SECRET", "APP_DEBUG"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in report output", key)
		}
	}
}
