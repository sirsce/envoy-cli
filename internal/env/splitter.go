package env

import "fmt"

// SplitStrategy controls how entries are split into buckets.
type SplitStrategy int

const (
	SplitByPrefix SplitStrategy = iota
	SplitByFirstChar
	SplitEven
)

// SplitResult holds a named bucket of entries.
type SplitResult struct {
	Name    string
	Entries []Entry
}

// Splitter partitions a slice of entries into named buckets.
type Splitter struct {
	strategy SplitStrategy
	buckets  int
}

// NewSplitter creates a Splitter with the given strategy.
// buckets is only used for SplitEven.
func NewSplitter(strategy SplitStrategy, buckets int) *Splitter {
	if buckets < 1 {
		buckets = 1
	}
	return &Splitter{strategy: strategy, buckets: buckets}
}

// Split divides entries according to the configured strategy.
func (s *Splitter) Split(entries []Entry) ([]SplitResult, error) {
	switch s.strategy {
	case SplitByPrefix:
		return s.splitByPrefix(entries), nil
	case SplitByFirstChar:
		return s.splitByFirstChar(entries), nil
	case SplitEven:
		return s.splitEven(entries), nil
	default:
		return nil, fmt.Errorf("splitter: unknown strategy %d", s.strategy)
	}
}

func (s *Splitter) splitByPrefix(entries []Entry) []SplitResult {
	index := map[string]int{}
	var results []SplitResult
	for _, e := range entries {
		prefix := groupPrefix(e.Key)
		if prefix == "" {
			prefix = "_"
		}
		idx, ok := index[prefix]
		if !ok {
			idx = len(results)
			index[prefix] = idx
			results = append(results, SplitResult{Name: prefix})
		}
		results[idx].Entries = append(results[idx].Entries, e)
	}
	return results
}

func (s *Splitter) splitByFirstChar(entries []Entry) []SplitResult {
	index := map[string]int{}
	var results []SplitResult
	for _, e := range entries {
		name := "_"
		if len(e.Key) > 0 {
			name = string(e.Key[0])
		}
		idx, ok := index[name]
		if !ok {
			idx = len(results)
			index[name] = idx
			results = append(results, SplitResult{Name: name})
		}
		results[idx].Entries = append(results[idx].Entries, e)
	}
	return results
}

func (s *Splitter) splitEven(entries []Entry) []SplitResult {
	results := make([]SplitResult, s.buckets)
	for i := range results {
		results[i].Name = fmt.Sprintf("bucket_%d", i+1)
	}
	for i, e := range entries {
		b := i % s.buckets
		results[b].Entries = append(results[b].Entries, e)
	}
	return results
}
