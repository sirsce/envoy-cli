package storage_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/storage"
)

func TestMemoryBackend_Empty(t *testing.T) {
	m := storage.NewMemoryBackend()

	exists, err := m.Exists()
	if err != nil {
		t.Fatalf("Exists error: %v", err)
	}
	if exists {
		t.Fatal("expected empty backend to not exist")
	}
}

func TestMemoryBackend_WriteRead(t *testing.T) {
	m := storage.NewMemoryBackend()

	const content = "DB_HOST=localhost\nDB_PORT=5432\n"
	if err := m.Write(strings.NewReader(content)); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	exists, err := m.Exists()
	if err != nil {
		t.Fatalf("Exists error: %v", err)
	}
	if !exists {
		t.Fatal("expected backend to exist after write")
	}

	if string(m.Bytes()) != content {
		t.Errorf("Bytes mismatch: got %q want %q", m.Bytes(), content)
	}

	rc, err := m.Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	defer rc.Close()

	buf := new(strings.Builder)
	if _, err := buf.ReadFrom(rc); err != nil {
		t.Fatalf("ReadFrom error: %v", err)
	}
	if buf.String() != content {
		t.Errorf("Read mismatch: got %q want %q", buf.String(), content)
	}
}

func TestMemoryBackend_OverwriteData(t *testing.T) {
	m := storage.NewMemoryBackend()

	_ = m.Write(strings.NewReader("OLD=data"))
	_ = m.Write(strings.NewReader("NEW=data"))

	if got := string(m.Bytes()); got != "NEW=data" {
		t.Errorf("expected overwrite, got %q", got)
	}
}
