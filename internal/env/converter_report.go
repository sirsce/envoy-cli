package env

import (
	"encoding/json"
	"fmt"
	"io"
)

// ConvertReport summarises the results of a conversion operation.
type ConvertReport struct {
	results []ConvertResult
	from    ConvertFormat
	to      ConvertFormat
}

// NewConvertReport creates a ConvertReport from a slice of ConvertResults.
func NewConvertReport(from, to ConvertFormat, results []ConvertResult) *ConvertReport {
	return &ConvertReport{results: results, from: from, to: to}
}

// SuccessCount returns the number of successfully converted entries.
func (r *ConvertReport) SuccessCount() int {
	count := 0
	for _, res := range r.results {
		if res.Success {
			count++
		}
	}
	return count
}

// WriteText writes a human-readable report to w.
func (r *ConvertReport) WriteText(w io.Writer) error {
	if len(r.results) == 0 {
		_, err := fmt.Fprintln(w, "no entries converted")
		return err
	}
	fmt.Fprintf(w, "converted %d entries from %s to %s\n", r.SuccessCount(), r.from, r.to)
	for _, res := range r.results {
		status := "ok"
		if !res.Success {
			status = fmt.Sprintf("error: %v", res.Err)
		}
		fmt.Fprintf(w, "  %-30s %s\n", res.Key, status)
	}
	return nil
}

// WriteJSON writes a JSON-encoded report to w.
func (r *ConvertReport) WriteJSON(w io.Writer) error {
	type jsonResult struct {
		Key     string `json:"key"`
		From    string `json:"from"`
		To      string `json:"to"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}
	payload := make([]jsonResult, 0, len(r.results))
	for _, res := range r.results {
		jr := jsonResult{
			Key:     res.Key,
			From:    string(res.From),
			To:      string(res.To),
			Success: res.Success,
		}
		if res.Err != nil {
			jr.Error = res.Err.Error()
		}
		payload = append(payload, jr)
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
