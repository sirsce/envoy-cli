package env

import (
	"fmt"
	"strings"
)

// ConvertFormat represents the target format for conversion.
type ConvertFormat string

const (
	FormatDotenv ConvertFormat = "dotenv"
	FormatExport ConvertFormat = "export"
	FormatJSON   ConvertFormat = "json"
	FormatYAML   ConvertFormat = "yaml"
)

// ConvertResult holds the result of a single entry conversion.
type ConvertResult struct {
	Key     string
	From    ConvertFormat
	To      ConvertFormat
	Success bool
	Err     error
}

// Converter transforms a slice of entries into a target format string.
type Converter struct {
	from ConvertFormat
	to   ConvertFormat
}

// NewConverter creates a Converter that converts from one format to another.
func NewConverter(from, to ConvertFormat) (*Converter, error) {
	valid := map[ConvertFormat]bool{
		FormatDotenv: true,
		FormatExport: true,
		FormatJSON:   true,
		FormatYAML:   true,
	}
	if !valid[from] {
		return nil, fmt.Errorf("unsupported source format: %s", from)
	}
	if !valid[to] {
		return nil, fmt.Errorf("unsupported target format: %s", to)
	}
	return &Converter{from: from, to: to}, nil
}

// Convert converts entries to the target format string and returns results.
func (c *Converter) Convert(entries []Entry) (string, []ConvertResult, error) {
	results := make([]ConvertResult, 0, len(entries))
	for _, e := range entries {
		results = append(results, ConvertResult{
			Key:     e.Key,
			From:    c.from,
			To:      c.to,
			Success: true,
		})
	}

	var sb strings.Builder
	switch c.to {
	case FormatDotenv:
		for _, e := range entries {
			fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
		}
	case FormatExport:
		for _, e := range entries {
			fmt.Fprintf(&sb, "export %s=%q\n", e.Key, e.Value)
		}
	case FormatJSON:
		sb.WriteString("{\n")
		for i, e := range entries {
			comma := ","
			if i == len(entries)-1 {
				comma = ""
			}
			fmt.Fprintf(&sb, "  %q: %q%s\n", e.Key, e.Value, comma)
		}
		sb.WriteString("}")
	case FormatYAML:
		for _, e := range entries {
			fmt.Fprintf(&sb, "%s: %q\n", e.Key, e.Value)
		}
	}

	return sb.String(), results, nil
}

// Formats returns the from and to formats of the converter.
func (c *Converter) Formats() (ConvertFormat, ConvertFormat) {
	return c.from, c.to
}
