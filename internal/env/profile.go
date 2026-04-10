package env

import "fmt"

// Profile represents a named set of environment entries (e.g., "development", "production").
type Profile struct {
	Name    string
	Entries []Entry
}

// ProfileManager manages multiple named profiles.
type ProfileManager struct {
	profiles map[string]*Profile
}

// NewProfileManager creates a new ProfileManager.
func NewProfileManager() *ProfileManager {
	return &ProfileManager{
		profiles: make(map[string]*Profile),
	}
}

// Add registers a profile under the given name.
func (pm *ProfileManager) Add(name string, entries []Entry) {
	pm.profiles[name] = &Profile{Name: name, Entries: entries}
}

// Get retrieves a profile by name.
func (pm *ProfileManager) Get(name string) (*Profile, error) {
	p, ok := pm.profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q not found", name)
	}
	return p, nil
}

// List returns all registered profile names.
func (pm *ProfileManager) List() []string {
	names := make([]string, 0, len(pm.profiles))
	for name := range pm.profiles {
		names = append(names, name)
	}
	return names
}

// Merge merges a base profile with an override profile.
// Keys present in override replace those in base.
func (pm *ProfileManager) Merge(baseName, overrideName string) ([]Entry, error) {
	base, err := pm.Get(baseName)
	if err != nil {
		return nil, err
	}
	override, err := pm.Get(overrideName)
	if err != nil {
		return nil, err
	}

	result := make(map[string]Entry)
	for _, e := range base.Entries {
		result[e.Key] = e
	}
	for _, e := range override.Entries {
		result[e.Key] = e
	}

	merged := make([]Entry, 0, len(result))
	for _, e := range result {
		merged = append(merged, e)
	}
	return merged, nil
}
