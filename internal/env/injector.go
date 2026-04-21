package env

import "fmt"

// InjectorStrategy controls how conflicts are handled during injection.
type InjectorStrategy int

const (
	// InjectorSkip skips entries that already exist in the target.
	InjectorSkip InjectorStrategy = iota
	// InjectorOverwrite overwrites entries that already exist in the target.
	InjectorOverwrite
	// InjectorError returns an error when a conflict is detected.
	InjectorError
)

// InjectorResult holds the outcome of a single injection operation.
type InjectorResult struct {
	Key      string
	Injected bool
	Skipped  bool
	Reason   string
}

// Injector merges a set of source entries into a target slice.
type Injector struct {
	strategy InjectorStrategy
}

// NewInjector creates a new Injector with the given strategy.
func NewInjector(strategy InjectorStrategy) *Injector {
	return &Injector{strategy: strategy}
}

// Inject injects source entries into target, returning the merged slice and results.
func (inj *Injector) Inject(target, source []Entry) ([]Entry, []InjectorResult, error) {
	index := make(map[string]int, len(target))
	for i, e := range target {
		index[e.Key] = i
	}

	out := make([]Entry, len(target))
	copy(out, target)

	var results []InjectorResult

	for _, src := range source {
		if idx, exists := index[src.Key]; exists {
			switch inj.strategy {
			case InjectorSkip:
				results = append(results, InjectorResult{Key: src.Key, Skipped: true, Reason: "key already exists"})
			case InjectorOverwrite:
				out[idx] = src
				results = append(results, InjectorResult{Key: src.Key, Injected: true, Reason: "overwritten"})
			case InjectorError:
				return nil, nil, fmt.Errorf("injector: conflict on key %q", src.Key)
			}
		} else {
			out = append(out, src)
			index[src.Key] = len(out) - 1
			results = append(results, InjectorResult{Key: src.Key, Injected: true, Reason: "new key"})
		}
	}

	return out, results, nil
}

// InjectedCount returns the number of successfully injected entries.
func InjectedCount(results []InjectorResult) int {
	n := 0
	for _, r := range results {
		if r.Injected {
			n++
		}
	}
	return n
}
