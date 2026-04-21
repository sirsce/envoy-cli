package env

import (
	"fmt"
	"io"
)

// ChainLoader loads and merges multiple .env sources in order,
// with later sources overriding earlier ones.
type ChainLoader struct {
	parser  *Parser
	sources []namedReader
	merger  *Merger
}

type namedReader struct {
	name   string
	reader io.Reader
}

// ChainResult holds the merged entries and metadata about each source.
type ChainResult struct {
	Entries []Entry
	Loaded  []string
	Skipped []string
}

// NewChainLoader creates a ChainLoader using the default parser and
// an Override merge strategy so later sources win on conflict.
func NewChainLoader() *ChainLoader {
	return &ChainLoader{
		parser: NewParser(),
		merger: NewMerger(MergeOverride),
	}
}

// Add registers a named reader as the next source in the chain.
func (cl *ChainLoader) Add(name string, r io.Reader) *ChainLoader {
	cl.sources = append(cl.sources, namedReader{name: name, reader: r})
	return cl
}

// Load processes all sources in registration order and returns
// the merged result. Sources that produce parse errors are skipped
// and recorded in ChainResult.Skipped.
func (cl *ChainLoader) Load() (*ChainResult, error) {
	result := &ChainResult{}
	base := []Entry{}

	for _, src := range cl.sources {
		entries, err := cl.parser.Parse(src.reader)
		if err != nil {
			result.Skipped = append(result.Skipped, fmt.Sprintf("%s: %v", src.name, err))
			continue
		}
		merged, err := cl.merger.Merge(base, entries)
		if err != nil {
			return nil, fmt.Errorf("chainloader: merging %q: %w", src.name, err)
		}
		base = merged
		result.Loaded = append(result.Loaded, src.name)
	}

	result.Entries = base
	return result, nil
}

// SourceCount returns the number of registered sources.
func (cl *ChainLoader) SourceCount() int {
	return len(cl.sources)
}
