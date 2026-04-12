package env

import (
	"fmt"
	"time"
)

// PinnedEntry represents an env entry that has been pinned to a specific value.
type PinnedEntry struct {
	Key       string
	Value     string
	PinnedAt  time.Time
	Reason    string
}

// Pinner manages pinned environment variable entries that should not be
// overwritten during sync or merge operations.
type Pinner struct {
	clock  func() time.Time
	pinned map[string]PinnedEntry
}

// PinnerOption is a functional option for Pinner.
type PinnerOption func(*Pinner)

// WithPinnerClock sets a custom clock for the Pinner (useful for testing).
func WithPinnerClock(fn func() time.Time) PinnerOption {
	return func(p *Pinner) {
		p.clock = fn
	}
}

// NewPinner creates a new Pinner instance.
func NewPinner(opts ...PinnerOption) *Pinner {
	p := &Pinner{
		clock:  time.Now,
		pinned: make(map[string]PinnedEntry),
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Pin pins a key to its current value with an optional reason.
func (p *Pinner) Pin(key, value, reason string) {
	p.pinned[key] = PinnedEntry{
		Key:      key,
		Value:    value,
		PinnedAt: p.clock(),
		Reason:   reason,
	}
}

// Unpin removes a pin for the given key.
func (p *Pinner) Unpin(key string) error {
	if _, ok := p.pinned[key]; !ok {
		return fmt.Errorf("pinner: key %q is not pinned", key)
	}
	delete(p.pinned, key)
	return nil
}

// IsPinned reports whether the given key is pinned.
func (p *Pinner) IsPinned(key string) bool {
	_, ok := p.pinned[key]
	return ok
}

// Get returns the PinnedEntry for a key, and whether it exists.
func (p *Pinner) Get(key string) (PinnedEntry, bool) {
	e, ok := p.pinned[key]
	return e, ok
}

// Apply takes a slice of entries and replaces any pinned keys with their
// pinned values, leaving non-pinned entries untouched.
func (p *Pinner) Apply(entries []Entry) []Entry {
	out := make([]Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if pe, ok := p.pinned[e.Key]; ok {
			out[i].Value = pe.Value
		}
	}
	return out
}

// List returns all currently pinned entries.
func (p *Pinner) List() []PinnedEntry {
	result := make([]PinnedEntry, 0, len(p.pinned))
	for _, e := range p.pinned {
		result = append(result, e)
	}
	return result
}
