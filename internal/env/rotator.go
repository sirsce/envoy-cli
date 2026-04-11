package env

import (
	"fmt"
	"time"
)

// RotationRecord captures a single key rotation event.
type RotationRecord struct {
	Key       string
	OldValue  string
	NewValue  string
	RotatedAt time.Time
}

// RotatorOption configures a Rotator.
type RotatorOption func(*Rotator)

// Rotator replaces values for specified keys using a generator function.
type Rotator struct {
	generators map[string]func(old string) (string, error)
	clock      func() time.Time
}

// NewRotator creates a Rotator with optional configuration.
func NewRotator(opts ...RotatorOption) *Rotator {
	r := &Rotator{
		generators: make(map[string]func(old string) (string, error)),
		clock:      time.Now,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// WithClock overrides the time source (useful for tests).
func WithClock(fn func() time.Time) RotatorOption {
	return func(r *Rotator) { r.clock = fn }
}

// Register associates a key with a generator function.
func (r *Rotator) Register(key string, gen func(old string) (string, error)) {
	r.generators[key] = gen
}

// Rotate applies registered generators to the provided entries.
// It returns updated entries and a log of all rotations performed.
func (r *Rotator) Rotate(entries []Entry) ([]Entry, []RotationRecord, error) {
	result := make([]Entry, len(entries))
	copy(result, entries)

	var records []RotationRecord

	for i, e := range result {
		gen, ok := r.generators[e.Key]
		if !ok {
			continue
		}
		newVal, err := gen(e.Value)
		if err != nil {
			return nil, nil, fmt.Errorf("rotator: key %q: %w", e.Key, err)
		}
		records = append(records, RotationRecord{
			Key:       e.Key,
			OldValue:  e.Value,
			NewValue:  newVal,
			RotatedAt: r.clock(),
		})
		result[i].Value = newVal
	}

	return result, records, nil
}

// Keys returns the list of registered rotation keys.
func (r *Rotator) Keys() []string {
	keys := make([]string, 0, len(r.generators))
	for k := range r.generators {
		keys = append(keys, k)
	}
	return keys
}
