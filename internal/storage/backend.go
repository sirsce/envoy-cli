package storage

import "io"

// Backend defines the interface for reading and writing .env data
// to a storage location (local file, remote S3, etc.).
type Backend interface {
	// Read returns a reader for the stored .env content.
	Read() (io.ReadCloser, error)

	// Write persists the .env content from the given reader.
	Write(r io.Reader) error

	// Exists reports whether the storage location already has data.
	Exists() (bool, error)
}
