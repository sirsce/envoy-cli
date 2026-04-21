package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// TokenReport summarises the output of a Tokenizer run.
type TokenReport struct {
	tokens []Token
}

// NewTokenReport wraps a slice of tokens for reporting.
func NewTokenReport(tokens []Token) *TokenReport {
	return &TokenReport{tokens: tokens}
}

// Count returns the total number of tokens.
func (r *TokenReport) Count() int { return len(r.tokens) }

// CountByType returns how many tokens match the given type.
func (r *TokenReport) CountByType(t TokenType) int {
	n := 0
	for _, tok := range r.tokens {
		if tok.Type == t {
			n++
		}
	}
	return n
}

// WriteText writes a human-readable summary to w.
func (r *TokenReport) WriteText(w io.Writer) error {
	if len(r.tokens) == 0 {
		_, err := fmt.Fprintln(w, "no tokens")
		return err
	}
	var sb strings.Builder
	for _, tok := range r.tokens {
		switch tok.Type {
		case TokenKey:
			sb.WriteString(fmt.Sprintf("[%d] KEY   %s = %s\n", tok.Line, tok.Key, tok.Value))
		case TokenComment:
			sb.WriteString(fmt.Sprintf("[%d] COMMENT %s\n", tok.Line, strings.TrimSpace(tok.Raw)))
		case TokenBlank:
			sb.WriteString(fmt.Sprintf("[%d] BLANK\n", tok.Line))
		}
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// WriteJSON writes a JSON array of token summaries to w.
func (r *TokenReport) WriteJSON(w io.Writer) error {
	type row struct {
		Line  int    `json:"line"`
		Type  string `json:"type"`
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}
	rows := make([]row, 0, len(r.tokens))
	for _, tok := range r.tokens {
		rows = append(rows, row{
			Line:  tok.Line,
			Type:  string(tok.Type),
			Key:   tok.Key,
			Value: tok.Value,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
