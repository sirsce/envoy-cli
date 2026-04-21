package env

import (
	"testing"
)

func TestTokenizer_KeyValue(t *testing.T) {
	tok := NewTokenizer()
	tokens, err := tok.Tokenize("FOO=bar\nBAZ=qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Key != "FOO" || tokens[0].Value != "bar" {
		t.Errorf("unexpected first token: %+v", tokens[0])
	}
	if tokens[1].Key != "BAZ" || tokens[1].Value != "qux" {
		t.Errorf("unexpected second token: %+v", tokens[1])
	}
}

func TestTokenizer_StripsQuotes(t *testing.T) {
	tok := NewTokenizer()
	tokens, err := tok.Tokenize(`SECRET="hello world"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens[0].Value != "hello world" {
		t.Errorf("expected unquoted value, got %q", tokens[0].Value)
	}
}

func TestTokenizer_SkipsCommentsAndBlanks(t *testing.T) {
	tok := NewTokenizer()
	tokens, err := tok.Tokenize("# comment\n\nFOO=1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(tokens))
	}
	if tokens[0].Key != "FOO" {
		t.Errorf("expected FOO, got %q", tokens[0].Key)
	}
}

func TestTokenizer_KeepComments(t *testing.T) {
	tok := NewTokenizer(WithKeepComments())
	tokens, err := tok.Tokenize("# a comment\nFOO=1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(tokens))
	}
	if tokens[0].Type != TokenComment {
		t.Errorf("expected COMMENT token, got %s", tokens[0].Type)
	}
}

func TestTokenizer_KeepBlanks(t *testing.T) {
	tok := NewTokenizer(WithKeepBlanks())
	tokens, err := tok.Tokenize("FOO=1\n\nBAR=2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tokens) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(tokens))
	}
	if tokens[1].Type != TokenBlank {
		t.Errorf("expected BLANK token, got %s", tokens[1].Type)
	}
}

func TestTokenizer_InvalidLine_MissingEquals(t *testing.T) {
	tok := NewTokenizer()
	_, err := tok.Tokenize("NOEQUALSSIGN")
	if err == nil {
		t.Fatal("expected error for missing '=', got nil")
	}
}

func TestTokenizer_LineNumbers(t *testing.T) {
	tok := NewTokenizer()
	tokens, err := tok.Tokenize("A=1\nB=2\nC=3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, tok := range tokens {
		if tok.Line != i+1 {
			t.Errorf("token %d: expected line %d, got %d", i, i+1, tok.Line)
		}
	}
}
