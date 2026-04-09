package env

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Writer handles writing .env files
type Writer struct{}

// NewWriter creates a new env file writer
func NewWriter() *Writer {
	return &Writer{}
}

// Write writes an EnvFile to the specified path
func (w *Writer) Write(filepath string, envFile *EnvFile) error {
	content := w.Format(envFile)

	err := os.WriteFile(filepath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Format converts an EnvFile to a formatted string
func (w *Writer) Format(envFile *EnvFile) string {
	var builder strings.Builder

	for _, entry := range envFile.Entries {
		// Quote value if it contains spaces
		value := entry.Value
		if strings.Contains(value, " ") {
			value = fmt.Sprintf(`"%s"`, value)
		}

		builder.WriteString(fmt.Sprintf("%s=%s\n", entry.Key, value))
	}

	return builder.String()
}

// WriteMap writes a map of key-value pairs to a .env file
func (w *Writer) WriteMap(filepath string, envMap map[string]string) error {
	// Convert map to EnvFile with sorted keys for consistent output
	keys := make([]string, 0, len(envMap))
	for k := range envMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{
			Key:   k,
			Value: envMap[k],
		})
	}

	envFile := &EnvFile{Entries: entries}
	return w.Write(filepath, envFile)
}

// Merge merges multiple EnvFiles into one, with later entries overwriting earlier ones
func Merge(envFiles ...*EnvFile) *EnvFile {
	merged := make(map[string]string)

	for _, envFile := range envFiles {
		for _, entry := range envFile.Entries {
			merged[entry.Key] = entry.Value
		}
	}

	// Convert back to entries with sorted keys
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]Entry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, Entry{
			Key:   k,
			Value: merged[k],
		})
	}

	return &EnvFile{Entries: entries}
}
