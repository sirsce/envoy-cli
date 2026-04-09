package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindConfigFile walks up the directory tree from start looking for
// a .envoy.json file, returning its path or an empty string if not found.
func FindConfigFile(start string) string {
	dir := start
	for {
		candidate := filepath.Join(dir, defaultConfigFile)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// LoadFromCWD attempts to load a config by searching from the current
// working directory upward. Falls back to DefaultConfig if not found.
func LoadFromCWD() (*Config, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, "", fmt.Errorf("config: cannot determine working directory: %w", err)
	}
	path := FindConfigFile(cwd)
	cfg, err := Load(path)
	if err != nil {
		return nil, path, fmt.Errorf("config: failed to load %q: %w", path, err)
	}
	return cfg, path, nil
}

// MustLoad loads a config from path and panics on error. Useful in tests.
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(fmt.Sprintf("config: MustLoad(%q): %v", path, err))
	}
	return cfg
}
