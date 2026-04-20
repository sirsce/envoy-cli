package env

// Pipeline chains multiple entry transformation steps together,
// executing them in order and passing the output of each step
// as the input to the next. This allows composing complex
// env-processing workflows from small, reusable components.

// PipelineStep is a function that transforms a slice of Entry values.
// It returns the transformed entries or an error.
type PipelineStep func(entries []Entry) ([]Entry, error)

// PipelineResult holds the final entries and per-step metadata
// collected during a pipeline run.
type PipelineResult struct {
	Entries []Entry
	Steps   []PipelineStepResult
}

// PipelineStepResult records the outcome of a single pipeline step.
type PipelineStepResult struct {
	Name    string
	Input   int
	Output  int
	Skipped bool
	Err     error
}

// Pipeline executes a sequence of PipelineSteps against an initial
// set of entries, collecting per-step results along the way.
type Pipeline struct {
	steps []namedStep
}

type namedStep struct {
	name string
	fn   PipelineStep
}

// NewPipeline creates an empty Pipeline ready to have steps added.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStep appends a named step to the pipeline.
// Steps are executed in the order they are added.
func (p *Pipeline) AddStep(name string, fn PipelineStep) *Pipeline {
	p.steps = append(p.steps, namedStep{name: name, fn: fn})
	return p
}

// Run executes all registered steps in order, starting with the
// provided entries. If any step returns an error the pipeline halts
// and returns the partial PipelineResult along with the error.
func (p *Pipeline) Run(entries []Entry) (PipelineResult, error) {
	result := PipelineResult{
		Entries: copyEntries(entries),
		Steps:   make([]PipelineStepResult, 0, len(p.steps)),
	}

	current := copyEntries(entries)

	for _, s := range p.steps {
		stepResult := PipelineStepResult{
			Name:  s.name,
			Input: len(current),
		}

		if s.fn == nil {
			stepResult.Skipped = true
			result.Steps = append(result.Steps, stepResult)
			continue
		}

		out, err := s.fn(current)
		if err != nil {
			stepResult.Err = err
			result.Steps = append(result.Steps, stepResult)
			result.Entries = current
			return result, err
		}

		stepResult.Output = len(out)
		result.Steps = append(result.Steps, stepResult)
		current = out
	}

	result.Entries = current
	return result, nil
}

// StepCount returns the number of steps registered in the pipeline.
func (p *Pipeline) StepCount() int {
	return len(p.steps)
}

// copyEntries returns a shallow copy of the entry slice so that
// pipeline steps cannot accidentally mutate shared state.
func copyEntries(src []Entry) []Entry {
	if src == nil {
		return nil
	}
	out := make([]Entry, len(src))
	copy(out, src)
	return out
}
