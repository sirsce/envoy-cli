package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".envoy.json"

// Config holds the envoy-cli project configuration.
type Config struct {
	Remote    string            `json:"remote"`
	Encrypt   bool              `json:"encrypt"`
	KeyFile   string            `json:"key_file,omitempty"`
	EnvFiles  []string          `json:"env_files"`
	Aliases   map[string]string `json:"aliases,omitempty"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Encrypt:  true,
		EnvFiles: []string{".env"},
		Aliases:  make(map[string]string),
	}
}

// Load reads a Config from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigFile
	}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return DefaultConfig(), nil
		}
		return nil, err
	}
	defer f.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the Config to the given path as JSON.
func Save(cfg *Config, path string) error {
	if path == "" {
		path = defaultConfigFile
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

// Validate checks that the Config contains required fields.
func (c *Config) Validate() error {
	if len(c.EnvFiles) == 0 {
		return errors.New("config: at least one env file must be specified")
	}
	if c.Encrypt && c.KeyFile == "" {
		return errors.New("config: key_file is required when encryption is enabled")
	}
	return nil
}
