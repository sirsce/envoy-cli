package env

import (
	"strings"
	"testing"
)

func TestParseReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Entry
		wantErr bool
	}{
		{
			name: "valid env file",
			input: `DATABASE_URL=postgres://localhost:5432/db
API_KEY=secret123
DEBUG=true`,
			want: []Entry{
				{Key: "DATABASE_URL", Value: "postgres://localhost:5432/db"},
				{Key: "API_KEY", Value: "secret123"},
				{Key: "DEBUG", Value: "true"},
			},
			wantErr: false,
		},
		{
			name: "with comments and empty lines",
			input: `# Database configuration
DATABASE_URL=postgres://localhost:5432/db

# API Settings
API_KEY=secret123`,
			want: []Entry{
				{Key: "DATABASE_URL", Value: "postgres://localhost:5432/db"},
				{Key: "API_KEY", Value: "secret123"},
			},
			wantErr: false,
		},
		{
			name:    "with quoted values",
			input:   `MESSAGE="Hello World"
PATH='/usr/local/bin'`,
			want: []Entry{
				{Key: "MESSAGE", Value: "Hello World"},
				{Key: "PATH", Value: "/usr/local/bin"},
			},
			wantErr: false,
		},
		{
			name:    "invalid format - no equals sign",
			input:   "INVALID_LINE",
			wantErr: true,
		},
		{
			name:    "invalid format - empty key",
			input:   "=value",
			wantErr: true,
		},
		{
			name:    "empty file",
			input:   "",
			want:    []Entry{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser()
			reader := strings.NewReader(tt.input)
			got, err := parser.ParseReader(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got.Entries) != len(tt.want) {
					t.Errorf("ParseReader() got %d entries, want %d", len(got.Entries), len(tt.want))
					return
				}

				for i, entry := range got.Entries {
					if entry.Key != tt.want[i].Key || entry.Value != tt.want[i].Value {
						t.Errorf("ParseReader() entry[%d] = %+v, want %+v", i, entry, tt.want[i])
					}
				}
			}
		})
	}
}

func TestToMap(t *testing.T) {
	envFile := &EnvFile{
		Entries: []Entry{
			{Key: "KEY1", Value: "value1"},
			{Key: "KEY2", Value: "value2"},
		},
	}

	got := envFile.ToMap()

	if len(got) != 2 {
		t.Errorf("ToMap() got %d entries, want 2", len(got))
	}

	if got["KEY1"] != "value1" {
		t.Errorf("ToMap()[KEY1] = %s, want value1", got["KEY1"])
	}

	if got["KEY2"] != "value2" {
		t.Errorf("ToMap()[KEY2] = %s, want value2", got["KEY2"])
	}
}
