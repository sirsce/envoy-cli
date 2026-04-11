package env

import (
	"fmt"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Interpolator expands variable references within env entry values.
type Interpolator struct {
	allowMissing bool
}

// InterpolatorOption configures an Interpolator.
type InterpolatorOption func(*Interpolator)

// WithAllowMissing prevents errors when a referenced variable is not found;
// the reference is left as-is instead.
func WithAllowMissing() InterpolatorOption {
	return func(i *Interpolator) {
		i.allowMissing = true
	}
}

// NewInterpolator creates a new Interpolator with the given options.
func NewInterpolator(opts ...InterpolatorOption) *Interpolator {
	i := &Interpolator{}
	for _, o := range opts {
		o(i)
	}
	return i
}

// Interpolate resolves variable references in entry values using the provided
// entries as the lookup table. It returns a new slice with expanded values.
func (i *Interpolator) Interpolate(entries []Entry) ([]Entry, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	result := make([]Entry, len(entries))
	for idx, e := range entries {
		expanded, err := i.expand(e.Value, lookup)
		if err != nil {
			return nil, fmt.Errorf("interpolate %q: %w", e.Key, err)
		}
		result[idx] = Entry{Key: e.Key, Value: expanded}
	}
	return result, nil
}

// expand replaces all variable references in s using the lookup map.
func (i *Interpolator) expand(s string, lookup map[string]string) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(s, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		val, ok := lookup[name]
		if !ok {
			if i.allowMissing {
				return match
			}
			expandErr = fmt.Errorf("undefined variable %q", name)
			return match
		}
		return val
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
