package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// PatchReport summarises the results of a Patcher.Apply call.
type PatchReport struct {
	Results []PatchResult
}

// NewPatchReport creates a PatchReport from a slice of PatchResult values.
func NewPatchReport(results []PatchResult) *PatchReport {
	return &PatchReport{Results: results}
}

// Applied returns only the results where the instruction was applied.
func (r *PatchReport) Applied() []PatchResult {
	var out []PatchResult
	for _, res := range r.Results {
		if res.Applied {
			out = append(out, res)
		}
	}
	return out
}

// Skipped returns results where the instruction was not applied.
func (r *PatchReport) Skipped() []PatchResult {
	var out []PatchResult
	for _, res := range r.Results {
		if !res.Applied {
			out = append(out, res)
		}
	}
	return out
}

// WriteText writes a human-readable summary to w.
func (r *PatchReport) WriteText(w io.Writer) error {
	if len(r.Results) == 0 {
		_, err := fmt.Fprintln(w, "No patch instructions applied.")
		return err
	}
	var sb strings.Builder
	for _, res := range r.Results {
		status := "OK"
		if !res.Applied {
			status = "SKIP"
		}
		sb.WriteString(fmt.Sprintf("[%s] %-8s %-20s %s\n",
			status, res.Instruction.Op, res.Instruction.Key, res.Note))
	}
	_, err := fmt.Fprint(w, sb.String())
	return err
}

// WriteJSON writes a JSON representation of the report to w.
func (r *PatchReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Op      string `json:"op"`
		Key     string `json:"key"`
		Applied bool   `json:"applied"`
		Note    string `json:"note"`
	}
	rows := make([]jsonResult, len(r.Results))
	for i, res := range r.Results {
		rows[i] = jsonResult{
			Op:      string(res.Instruction.Op),
			Key:     res.Instruction.Key,
			Applied: res.Applied,
			Note:    res.Note,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}
