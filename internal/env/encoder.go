package env

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncoderMode defines how values are encoded.
type EncoderMode int

const (
	EncodeBase64 EncoderMode = iota
	EncodeHex
	EncodeURL
)

// EncodeResult holds the result of encoding a single entry.
type EncodeResult struct {
	Key      string
	Original string
	Encoded  string
	Skipped  bool
}

// Encoder encodes entry values using a specified mode.
type Encoder struct {
	mode    EncoderMode
	keys    map[string]struct{}
	results []EncodeResult
}

// NewEncoder creates an Encoder targeting specific keys.
// If no keys are provided, all entries are encoded.
func NewEncoder(mode EncoderMode, keys ...string) *Encoder {
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[k] = struct{}{}
	}
	return &Encoder{mode: mode, keys: km}
}

// Encode encodes matching entry values and returns new entries.
func (e *Encoder) Encode(entries []Entry) ([]Entry, error) {
	e.results = nil
	out := make([]Entry, len(entries))
	for i, en := range entries {
		if len(e.keys) > 0 {
			if _, ok := e.keys[en.Key]; !ok {
				out[i] = en
				e.results = append(e.results, EncodeResult{Key: en.Key, Original: en.Value, Encoded: en.Value, Skipped: true})
				continue
			}
		}
		encoded, err := e.encodeValue(en.Value)
		if err != nil {
			return nil, fmt.Errorf("encoder: key %q: %w", en.Key, err)
		}
		e.results = append(e.results, EncodeResult{Key: en.Key, Original: en.Value, Encoded: encoded})
		out[i] = Entry{Key: en.Key, Value: encoded}
	}
	return out, nil
}

// Results returns the encode results from the last Encode call.
func (e *Encoder) Results() []EncodeResult {
	return e.results
}

// EncodedCount returns the number of entries that were encoded.
func (e *Encoder) EncodedCount() int {
	count := 0
	for _, r := range e.results {
		if !r.Skipped {
			count++
		}
	}
	return count
}

func (e *Encoder) encodeValue(v string) (string, error) {
	switch e.mode {
	case EncodeBase64:
		return base64.StdEncoding.EncodeToString([]byte(v)), nil
	case EncodeHex:
		var sb strings.Builder
		for _, b := range []byte(v) {
			fmt.Fprintf(&sb, "%02x", b)
		}
		return sb.String(), nil
	case EncodeURL:
		return base64.URLEncoding.EncodeToString([]byte(v)), nil
	default:
		return "", fmt.Errorf("unknown encoder mode: %d", e.mode)
	}
}
