package env

import (
	"fmt"
	"sort"
)

// ProfileRegistry combines a ProfileLoader with an in-memory ProfileManager,
// providing a unified interface for loading, caching, and merging profiles.
type ProfileRegistry struct {
	loader  *ProfileLoader
	manager *ProfileManager
}

// NewProfileRegistry creates a ProfileRegistry backed by the given directory.
func NewProfileRegistry(dir string) *ProfileRegistry {
	return &ProfileRegistry{
		loader:  NewProfileLoader(dir),
		manager: NewProfileManager(),
	}
}

// Ensure loads a profile from disk if not already cached.
func (pr *ProfileRegistry) Ensure(name string) error {
	if _, err := pr.manager.Get(name); err == nil {
		return nil // already cached
	}
	p, err := pr.loader.Load(name)
	if err != nil {
		return fmt.Errorf("registry: load profile %q: %w", name, err)
	}
	pr.manager.Add(p.Name, p.Entries)
	return nil
}

// Get returns a cached profile, loading it from disk if necessary.
func (pr *ProfileRegistry) Get(name string) (*Profile, error) {
	if err := pr.Ensure(name); err != nil {
		return nil, err
	}
	return pr.manager.Get(name)
}

// MergeProfiles merges two profiles (loading from disk if needed) and returns
// the combined entries with override taking precedence over base.
func (pr *ProfileRegistry) MergeProfiles(baseName, overrideName string) ([]Entry, error) {
	for _, name := range []string{baseName, overrideName} {
		if err := pr.Ensure(name); err != nil {
			return nil, err
		}
	}
	return pr.manager.Merge(baseName, overrideName)
}

// Names returns all currently cached profile names in sorted order.
func (pr *ProfileRegistry) Names() []string {
	names := pr.manager.List()
	sort.Strings(names)
	return names
}
