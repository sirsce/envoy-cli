package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalBackend stores .env data on the local filesystem.
type LocalBackend struct {
	path string
}

// NewLocalBackend creates a LocalBackend for the given file path.
func NewLocalBackend(path string) *LocalBackend {
	return &LocalBackend{path: filepath.Clean(path)}
}

// Read opens the local file and returns a ReadCloser.
func (l *LocalBackend) Read() (io.ReadCloser, error) {
	f, err := os.Open(l.path)
	if err != nil {
		return nil, fmt.Errorf("local read: %w", err)
	}
	return f, nil
}

// Write creates or truncates the local file and writes all data from r.
func (l *LocalBackend) Write(r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return fmt.Errorf("local write mkdir: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("local write open: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("local write copy: %w", err)
	}
	return nil
}

// Exists reports whether the file exists on disk.
func (l *LocalBackend) Exists() (bool, error) {
	_, err := os.Stat(l.path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("local exists: %w", err)
	}
	return true, nil
}
