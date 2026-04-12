package env

import "fmt"

// Scope represents a named environment scope (e.g. "development", "staging", "production").
type Scope struct {
	Name    string
	Entries []Entry
}

// ScopeManager manages multiple named scopes and allows resolving entries
// with optional fallback to a base scope.
type ScopeManager struct {
	scopes map[string]*Scope
	base   string
}

// NewScopeManager creates a new ScopeManager with an optional base scope name.
// The base scope is used as a fallback when resolving entries.
func NewScopeManager(base string) *ScopeManager {
	return &ScopeManager{
		scopes: make(map[string]*Scope),
		base:   base,
	}
}

// Add registers a named scope with the given entries.
func (sm *ScopeManager) Add(name string, entries []Entry) {
	sm.scopes[name] = &Scope{Name: name, Entries: entries}
}

// Get returns the entries for the named scope, falling back to the base scope
// for any keys not present in the requested scope. Returns an error if the
// requested scope does not exist.
func (sm *ScopeManager) Get(name string) ([]Entry, error) {
	scope, ok := sm.scopes[name]
	if !ok {
		return nil, fmt.Errorf("scope %q not found", name)
	}

	if sm.base == "" || name == sm.base {
		return scope.Entries, nil
	}

	base, hasBase := sm.scopes[sm.base]
	if !hasBase {
		return scope.Entries, nil
	}

	// Build a map from the requested scope for fast lookup.
	resolved := make(map[string]Entry, len(scope.Entries))
	for _, e := range scope.Entries {
		resolved[e.Key] = e
	}

	// Fill in missing keys from the base scope.
	for _, e := range base.Entries {
		if _, exists := resolved[e.Key]; !exists {
			resolved[e.Key] = e
		}
	}

	// Convert back to a slice preserving base order for new keys.
	seen := make(map[string]bool)
	var result []Entry
	for _, e := range scope.Entries {
		result = append(result, resolved[e.Key])
		seen[e.Key] = true
	}
	for _, e := range base.Entries {
		if !seen[e.Key] {
			result = append(result, resolved[e.Key])
		}
	}
	return result, nil
}

// List returns all registered scope names.
func (sm *ScopeManager) List() []string {
	names := make([]string, 0, len(sm.scopes))
	for name := range sm.scopes {
		names = append(names, name)
	}
	return names
}
