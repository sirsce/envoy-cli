package env

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkHashFile measures the cost of hashing a typical .env file.
func BenchmarkHashFile(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, ".env")

	// Build a realistic 50-entry .env file.
	var content string
	for i := 0; i < 50; i++ {
		content += fmt.Sprintf("ENV_VAR_%02d=some_value_%02d\n", i, i)
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		b.Fatalf("WriteFile: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := hashFile(path); err != nil {
			b.Fatalf("hashFile: %v", err)
		}
	}
}

// BenchmarkWatcher_Seed measures the Seed path (hash + store).
func BenchmarkWatcher_Seed(b *testing.B) {
	dir := b.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte("KEY=value\n"), 0600); err != nil {
		b.Fatalf("WriteFile: %v", err)
	}

	w := NewWatcher(path, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := w.Seed(); err != nil {
			b.Fatalf("Seed: %v", err)
		}
	}
}
