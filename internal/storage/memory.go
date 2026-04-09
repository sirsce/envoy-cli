package storage

import (
	"bytes"
	"io"
)

// MemoryBackend is an in-memory Backend useful for testing.
type MemoryBackend struct {
	buf []byte
}

// NewMemoryBackend returns an empty MemoryBackend.
func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{}
}

// Read returns a reader over the stored bytes.
func (m *MemoryBackend) Read() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(m.buf)), nil
}

// Write replaces the stored bytes with the content from r.
func (m *MemoryBackend) Write(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	m.buf = data
	return nil
}

// Exists reports whether any data has been written.
func (m *MemoryBackend) Exists() (bool, error) {
	return len(m.buf) > 0, nil
}

// Bytes returns a copy of the stored bytes (useful in tests).
func (m *MemoryBackend) Bytes() []byte {
	out := make([]byte, len(m.buf))
	copy(out, m.buf)
	return out
}
