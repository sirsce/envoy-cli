package env

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ProfileLoader loads profiles from .env files on disk.
// It looks for files named <name>.env or .<name>.env in the given directory.
type ProfileLoader struct {
	dir    string
	parser *Parser
}

// NewProfileLoader creates a ProfileLoader rooted at dir.
func NewProfileLoader(dir string) *ProfileLoader {
	return &ProfileLoader{dir: dir, parser: NewParser()}
}

// Load reads a profile by name from the filesystem.
func (pl *ProfileLoader) Load(name string) (*Profile, error) {
	path, err := pl.resolvePath(name)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open profile file: %w", err)
	}
	defer f.Close()

	entries, err := pl.parseEntries(f)
	if err != nil {
		return nil, err
	}
	return &Profile{Name: name, Entries: entries}, nil
}

func (pl *ProfileLoader) resolvePath(name string) (string, error) {
	candidates := []string{
		filepath.Join(pl.dir, name+".env"),
		filepath.Join(pl.dir, "."+name+".env"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("no env file found for profile %q in %s", name, pl.dir)
}

func (pl *ProfileLoader) parseEntries(r io.Reader) ([]Entry, error) {
	entries, err := pl.parser.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse profile: %w", err)
	}
	return entries, nil
}
