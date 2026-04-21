package env

import (
	"fmt"
	"strconv"
	"strings"
)

// CastType represents the target type for casting.
type CastType string

const (
	CastString  CastType = "string"
	CastInt     CastType = "int"
	CastFloat   CastType = "float"
	CastBool    CastType = "bool"
)

// CastRule defines a key and the type it should be cast to.
type CastRule struct {
	Key      string
	CastTo   CastType
}

// CastResult holds the result of a single cast operation.
type CastResult struct {
	Key     string
	From    string
	To      string
	Type    CastType
	Skipped bool
	Err     error
}

// Caster normalises env entry values to canonical string representations
// of a given type (e.g. "1" → "true" for bool, " 42 " → "42" for int).
type Caster struct {
	rules []CastRule
}

// NewCaster creates a Caster with the provided rules.
func NewCaster(rules []CastRule) *Caster {
	return &Caster{rules: rules}
}

// Apply iterates over entries and applies cast rules by key.
// Entries without a matching rule are returned unchanged.
func (c *Caster) Apply(entries []Entry) ([]Entry, []CastResult) {
	ruleMap := make(map[string]CastType, len(c.rules))
	for _, r := range c.rules {
		ruleMap[r.Key] = r.CastTo
	}

	out := make([]Entry, len(entries))
	copy(out, entries)

	var results []CastResult
	for i, e := range out {
		typ, ok := ruleMap[e.Key]
		if !ok {
			continue
		}
		norm, err := castValue(e.Value, typ)
		res := CastResult{Key: e.Key, From: e.Value, Type: typ}
		if err != nil {
			res.Err = err
			res.Skipped = true
		} else {
			res.To = norm
			out[i].Value = norm
		}
		results = append(results, res)
	}
	return out, results
}

// CastedCount returns the number of successfully cast entries from results.
func CastedCount(results []CastResult) int {
	n := 0
	for _, r := range results {
		if !r.Skipped {
			n++
		}
	}
	return n
}

func castValue(v string, typ CastType) (string, error) {
	v = strings.TrimSpace(v)
	switch typ {
	case CastString:
		return v, nil
	case CastInt:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int: %w", v, err)
		}
		return strconv.FormatInt(n, 10), nil
	case CastFloat:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float: %w", v, err)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CastBool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to bool: %w", v, err)
		}
		if b {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("unknown cast type %q", typ)
	}
}
