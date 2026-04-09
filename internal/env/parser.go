package env

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Entry represents a single environment variable entry
type Entry struct {
	Key   string
	Value string
}

// EnvFile represents a parsed .env file
type EnvFile struct {
	Entries []Entry
}

// Parser handles parsing of .env files
type Parser struct{}

// NewParser creates a new env file parser
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads and parses an .env file from the given path
func (p *Parser) Parse(filepath string) (*EnvFile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return p.ParseReader(file)
}

// ParseReader parses .env content from an io.Reader
func (p *Parser) ParseReader(r io.Reader) (*EnvFile, error) {
	envFile := &EnvFile{
		Entries: make([]Entry, 0),
	}

	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid format at line %d: %s", lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, `"`)
		value = strings.Trim(value, `'`)

		if key == "" {
			return nil, fmt.Errorf("empty key at line %d", lineNum)
		}

		envFile.Entries = append(envFile.Entries, Entry{
			Key:   key,
			Value: value,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return envFile, nil
}

// ToMap converts the EnvFile entries to a map
func (e *EnvFile) ToMap() map[string]string {
	result := make(map[string]string)
	for _, entry := range e.Entries {
		result[entry.Key] = entry.Value
	}
	return result
}
