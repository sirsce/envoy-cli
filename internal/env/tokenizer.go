package env

import (
	"fmt"
	"strings"
)

// TokenType represents the kind of a parsed token.
type TokenType string

const (
	TokenKey     TokenType = "KEY"
	TokenValue   TokenType = "VALUE"
	TokenComment TokenType = "COMMENT"
	TokenBlank   TokenType = "BLANK"
)

// Token is a single lexical unit from a .env file line.
type Token struct {
	Type    TokenType
	Raw     string
	Key     string
	Value   string
	Line    int
}

// Tokenizer breaks raw .env content into tokens.
type Tokenizer struct {
	keepComments bool
	keepBlanks   bool
}

// TokenizerOption configures a Tokenizer.
type TokenizerOption func(*Tokenizer)

// WithKeepComments instructs the tokenizer to emit comment tokens.
func WithKeepComments() TokenizerOption {
	return func(t *Tokenizer) { t.keepComments = true }
}

// WithKeepBlanks instructs the tokenizer to emit blank-line tokens.
func WithKeepBlanks() TokenizerOption {
	return func(t *Tokenizer) { t.keepBlanks = true }
}

// NewTokenizer creates a Tokenizer with optional configuration.
func NewTokenizer(opts ...TokenizerOption) *Tokenizer {
	t := &Tokenizer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Tokenize parses the given source string and returns a slice of Tokens.
func (t *Tokenizer) Tokenize(source string) ([]Token, error) {
	lines := strings.Split(source, "\n")
	var tokens []Token
	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			if t.keepBlanks {
				tokens = append(tokens, Token{Type: TokenBlank, Raw: line, Line: lineNum})
			}
		case strings.HasPrefix(trimmed, "#"):
			if t.keepComments {
				tokens = append(tokens, Token{Type: TokenComment, Raw: line, Line: lineNum})
			}
		default:
			tok, err := parseKeyValueToken(line, lineNum)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNum, err)
			}
			tokens = append(tokens, tok)
		}
	}
	return tokens, nil
}

func parseKeyValueToken(line string, lineNum int) (Token, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return Token{}, fmt.Errorf("missing '=' separator")
	}
	key := strings.TrimSpace(line[:idx])
	val := strings.TrimSpace(line[idx+1:])
	if key == "" {
		return Token{}, fmt.Errorf("empty key")
	}
	// Strip surrounding quotes from value.
	if len(val) >= 2 {
		if (val[0] == '"' && val[len(val)-1] == '"') ||
			(val[0] == '\'' && val[len(val)-1] == '\'') {
			val = val[1 : len(val)-1]
		}
	}
	return Token{Type: TokenKey, Raw: line, Key: key, Value: val, Line: lineNum}, nil
}
