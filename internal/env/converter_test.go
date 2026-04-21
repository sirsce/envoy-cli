package env

import (
	"strings"
	"testing"
)

func makeConverterEntries() []Entry {
	return []Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestConverter_InvalidFormats(t *testing.T) {
	_, err := NewConverter("bad", FormatJSON)
	if err == nil {
		t.Fatal("expected error for invalid source format")
	}
	_, err = NewConverter(FormatDotenv, "bad")
	if err == nil {
		t.Fatal("expected error for invalid target format")
	}
}

func TestConverter_ToDotenv(t *testing.T) {
	c, err := NewConverter(FormatDotenv, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, results, err := c.Convert(makeConverterEntries())
	if err != nil {
		t.Fatalf("convert error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line in output, got: %s", out)
	}
}

func TestConverter_ToExport(t *testing.T) {
	c, _ := NewConverter(FormatDotenv, FormatExport)
	out, _, _ := c.Convert(makeConverterEntries())
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected export prefix, got: %s", out)
	}
}

func TestConverter_ToJSON(t *testing.T) {
	c, _ := NewConverter(FormatDotenv, FormatJSON)
	out, _, _ := c.Convert(makeConverterEntries())
	if !strings.Contains(out, "\"APP_ENV\"") {
		t.Errorf("expected JSON key, got: %s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestConverter_ToYAML(t *testing.T) {
	c, _ := NewConverter(FormatDotenv, FormatYAML)
	out, _, _ := c.Convert(makeConverterEntries())
	if !strings.Contains(out, "APP_ENV:") {
		t.Errorf("expected YAML key, got: %s", out)
	}
}

func TestConverter_Formats(t *testing.T) {
	c, _ := NewConverter(FormatDotenv, FormatJSON)
	from, to := c.Formats()
	if from != FormatDotenv || to != FormatJSON {
		t.Errorf("unexpected formats: from=%s to=%s", from, to)
	}
}

func TestConvertReport_WriteText(t *testing.T) {
	c, _ := NewConverter(FormatDotenv, FormatJSON)
	_, results, _ := c.Convert(makeConverterEntries())
	report := NewConvertReport(FormatDotenv, FormatJSON, results)
	if report.SuccessCount() != 3 {
		t.Errorf("expected 3 successes, got %d", report.SuccessCount())
	}
	var sb strings.Builder
	if err := report.WriteText(&sb); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(sb.String(), "converted 3") {
		t.Errorf("unexpected report text: %s", sb.String())
	}
}

func TestConvertReport_WriteJSON(t *testing.T) {
	c, _ := NewConverter(FormatDotenv, FormatYAML)
	_, results, _ := c.Convert(makeConverterEntries())
	report := NewConvertReport(FormatDotenv, FormatYAML, results)
	var sb strings.Builder
	if err := report.WriteJSON(&sb); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	if !strings.Contains(sb.String(), "\"key\"") {
		t.Errorf("expected JSON output, got: %s", sb.String())
	}
}

func TestConvertReport_Empty(t *testing.T) {
	report := NewConvertReport(FormatDotenv, FormatJSON, nil)
	var sb strings.Builder
	if err := report.WriteText(&sb); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	if !strings.Contains(sb.String(), "no entries") {
		t.Errorf("expected empty message, got: %s", sb.String())
	}
}
