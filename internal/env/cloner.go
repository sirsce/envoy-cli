package env

import "fmt"

// CloneResult holds the outcome of a single key clone operation.
type CloneResult struct {
	SourceKey string
	DestKey   string
	Skipped   bool
	Reason    string
}

// Cloner duplicates env entries under new keys, optionally overwriting.
type Cloner struct {
	entries    []Entry
	overwrite  bool
	results    []CloneResult
}

// NewCloner creates a Cloner operating on the given entries.
func NewCloner(entries []Entry, overwrite bool) *Cloner {
	copied := make([]Entry, len(entries))
	copy(copied, entries)
	return &Cloner{entries: copied, overwrite: overwrite}
}

// Clone duplicates the value of srcKey into destKey.
// Returns an error if srcKey does not exist.
func (c *Cloner) Clone(srcKey, destKey string) error {
	if srcKey == "" || destKey == "" {
		return fmt.Errorf("cloner: source and destination keys must not be empty")
	}

	var srcValue string
	found := false
	for _, e := range c.entries {
		if e.Key == srcKey {
			srcValue = e.Value
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("cloner: source key %q not found", srcKey)
	}

	for i, e := range c.entries {
		if e.Key == destKey {
			if !c.overwrite {
				c.results = append(c.results, CloneResult{
					SourceKey: srcKey,
					DestKey:   destKey,
					Skipped:   true,
					Reason:    "destination already exists",
				})
				return nil
			}
			c.entries[i].Value = srcValue
			c.results = append(c.results, CloneResult{SourceKey: srcKey, DestKey: destKey})
			return nil
		}
	}

	c.entries = append(c.entries, Entry{Key: destKey, Value: srcValue})
	c.results = append(c.results, CloneResult{SourceKey: srcKey, DestKey: destKey})
	return nil
}

// Entries returns the current (possibly modified) entries.
func (c *Cloner) Entries() []Entry {
	out := make([]Entry, len(c.entries))
	copy(out, c.entries)
	return out
}

// Results returns the clone operation results recorded so far.
func (c *Cloner) Results() []CloneResult {
	return c.results
}

// AppliedCount returns the number of non-skipped clone operations.
func (c *Cloner) AppliedCount() int {
	count := 0
	for _, r := range c.results {
		if !r.Skipped {
			count++
		}
	}
	return count
}
