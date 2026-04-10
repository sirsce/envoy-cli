package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ExportFormat defines the output format for exported env entries.
type ExportFormat int

const (
	// FormatDotenv exports in standard KEY=VALUE format.
	FormatDotenv ExportFormat = iota
	// FormatExport exports with shell `export` prefix.
	FormatExport
	// FormatJSON exports as a JSON object.
	FormatJSON
)

// Exporter writes env entries to an io.Writer in a specified format.
type Exporter struct {
	format ExportFormat
	sorted bool
}

// NewExporter creates a new Exporter with the given format.
// If sorted is true, keys are written in alphabetical order.
func NewExporter(format ExportFormat, sorted bool) *Exporter {
	return &Exporter{format: format, sorted: sorted}
}

// Export writes the provided entries to w in the configured format.
func (e *Exporter) Export(w io.Writer, entries []Entry) error {
	switch e.format {
	case FormatDotenv:
		return e.writeDotenv(w, entries)
	case FormatExport:
		return e.writeExport(w, entries)
	case FormatJSON:
		return e.writeJSON(w, entries)
	default:
		return fmt.Errorf("unknown export format: %d", e.format)
	}
}

func (e *Exporter) ordered(entries []Entry) []Entry {
	if !e.sorted {
		return entries
	}
	copy := make([]Entry, len(entries))
	for i, en := range entries {
		copy[i] = en
	}
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Key < copy[j].Key
	})
	return copy
}

func (e *Exporter) writeDotenv(w io.Writer, entries []Entry) error {
	for _, en := range e.ordered(entries) {
		if _, err := fmt.Fprintf(w, "%s=%s\n", en.Key, en.Value); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeExport(w io.Writer, entries []Entry) error {
	for _, en := range e.ordered(entries) {
		if _, err := fmt.Fprintf(w, "export %s=%s\n", en.Key, en.Value); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(w io.Writer, entries []Entry) error {
	ordered := e.ordered(entries)
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, en := range ordered {
		value := strings.ReplaceAll(en.Value, `"`, `\"`)
		sb.WriteString(fmt.Sprintf("  \"%s\": \"%s\"", en.Key, value))
		if i < len(ordered)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}
