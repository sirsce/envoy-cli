package env

import (
	"fmt"
	"io"
	"strings"
)

// ReportFormat specifies the output format for lint reports.
type ReportFormat string

const (
	ReportFormatText ReportFormat = "text"
	ReportFormatJSON ReportFormat = "json"
)

// LintReport renders lint results to a writer.
type LintReport struct {
	format ReportFormat
}

// NewLintReport creates a LintReport with the given format.
func NewLintReport(format ReportFormat) *LintReport {
	return &LintReport{format: format}
}

// Write outputs the lint results to w.
func (r *LintReport) Write(w io.Writer, results []LintResult) error {
	switch r.format {
	case ReportFormatJSON:
		return r.writeJSON(w, results)
	default:
		return r.writeText(w, results)
	}
}

func (r *LintReport) writeText(w io.Writer, results []LintResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No lint violations found.")
		return err
	}
	for _, res := range results {
		_, err := fmt.Fprintf(w, "[%s] %s\n", res.Rule, res.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *LintReport) writeJSON(w io.Writer, results []LintResult) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, res := range results {
		sb.WriteString(fmt.Sprintf(
			"  {\"rule\": %q, \"key\": %q, \"message\": %q}",
			res.Rule, res.Entry.Key, res.Message,
		))
		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(w, sb.String())
	return err
}
