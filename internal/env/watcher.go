package env

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// WatchEvent describes a change detected in a watched file.
type WatchEvent struct {
	Path    string
	Changed bool
	Err     error
}

// Watcher polls a file for changes based on its content hash.
type Watcher struct {
	mu       sync.Mutex
	path     string
	lastHash string
	interval time.Duration
	stopCh   chan struct{}
}

// NewWatcher creates a Watcher for the given file path and poll interval.
func NewWatcher(path string, interval time.Duration) *Watcher {
	return &Watcher{
		path:     path,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins polling and sends events to the returned channel.
// The caller must call Stop to release resources.
func (w *Watcher) Start() <-chan WatchEvent {
	ch := make(chan WatchEvent, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stopCh:
				return
			case <-ticker.C:
				event := w.check()
				if event.Changed || event.Err != nil {
					ch <- event
				}
			}
		}
	}()
	return ch
}

// Stop halts the watcher goroutine.
func (w *Watcher) Stop() {
	close(w.stopCh)
}

// check reads the file and compares its hash to the last known hash.
func (w *Watcher) check() WatchEvent {
	hash, err := hashFile(w.path)
	if err != nil {
		return WatchEvent{Path: w.path, Err: err}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	changed := hash != w.lastHash
	w.lastHash = hash
	return WatchEvent{Path: w.path, Changed: changed}
}

// Seed initialises the baseline hash without emitting an event.
func (w *Watcher) Seed() error {
	hash, err := hashFile(w.path)
	if err != nil {
		return err
	}
	w.mu.Lock()
	w.lastHash = hash
	w.mu.Unlock()
	return nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
