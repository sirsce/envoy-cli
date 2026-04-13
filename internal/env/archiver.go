package env

import (
	"fmt"
	"sort"
	"time"
)

// ArchiveEntry represents a versioned snapshot of env entries.
type ArchiveEntry struct {
	Version   int
	Label     string
	Entries   []Entry
	ArchivedAt time.Time
}

// Archiver stores and retrieves versioned archives of env entry sets.
type Archiver struct {
	clock    func() time.Time
	archives []ArchiveEntry
	nextVer  int
}

// ArchiverOption configures an Archiver.
type ArchiverOption func(*Archiver)

// WithArchiverClock overrides the clock used for timestamps.
func WithArchiverClock(fn func() time.Time) ArchiverOption {
	return func(a *Archiver) { a.clock = fn }
}

// NewArchiver creates a new Archiver.
func NewArchiver(opts ...ArchiverOption) *Archiver {
	a := &Archiver{
		clock:   time.Now,
		nextVer: 1,
	}
	for _, o := range opts {
		o(a)
	}
	return a
}

// Archive stores a labeled snapshot of entries and returns the assigned version.
func (a *Archiver) Archive(label string, entries []Entry) int {
	copied := make([]Entry, len(entries))
	copy(copied, entries)
	v := a.nextVer
	a.archives = append(a.archives, ArchiveEntry{
		Version:    v,
		Label:      label,
		Entries:    copied,
		ArchivedAt: a.clock(),
	})
	a.nextVer++
	return v
}

// Get returns the ArchiveEntry for the given version, or an error if not found.
func (a *Archiver) Get(version int) (ArchiveEntry, error) {
	for _, ae := range a.archives {
		if ae.Version == version {
			return ae, nil
		}
	}
	return ArchiveEntry{}, fmt.Errorf("archive version %d not found", version)
}

// List returns all archived entries sorted by version ascending.
func (a *Archiver) List() []ArchiveEntry {
	out := make([]ArchiveEntry, len(a.archives))
	copy(out, a.archives)
	sort.Slice(out, func(i, j int) bool { return out[i].Version < out[j].Version })
	return out
}

// Delete removes the archive at the given version. Returns error if not found.
func (a *Archiver) Delete(version int) error {
	for i, ae := range a.archives {
		if ae.Version == version {
			a.archives = append(a.archives[:i], a.archives[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("archive version %d not found", version)
}

// Latest returns the most recently created ArchiveEntry, or an error if empty.
func (a *Archiver) Latest() (ArchiveEntry, error) {
	if len(a.archives) == 0 {
		return ArchiveEntry{}, fmt.Errorf("no archives exist")
	}
	latest := a.archives[0]
	for _, ae := range a.archives[1:] {
		if ae.Version > latest.Version {
			latest = ae
		}
	}
	return latest, nil
}
