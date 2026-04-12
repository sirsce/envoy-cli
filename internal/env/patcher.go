package env

import "fmt"

// PatchOp represents a single patch operation type.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchInstruction describes a single mutation to apply to an env entry set.
type PatchInstruction struct {
	Op      PatchOp
	Key     string
	Value   string // used for Set
	NewKey  string // used for Rename
}

// PatchResult records the outcome of applying a single instruction.
type PatchResult struct {
	Instruction PatchInstruction
	Applied     bool
	Note        string
}

// Patcher applies a sequence of patch instructions to a slice of Entry values.
type Patcher struct {
	instructions []PatchInstruction
}

// NewPatcher creates a Patcher with the given instructions.
func NewPatcher(instructions []PatchInstruction) *Patcher {
	return &Patcher{instructions: instructions}
}

// Apply executes all instructions against entries and returns the mutated
// slice along with a result record for each instruction.
func (p *Patcher) Apply(entries []Entry) ([]Entry, []PatchResult, error) {
	results := make([]PatchResult, 0, len(p.instructions))
	working := make([]Entry, len(entries))
	copy(working, entries)

	for _, inst := range p.instructions {
		var res PatchResult
		res.Instruction = inst

		switch inst.Op {
		case PatchSet:
			working, res = applySet(working, inst)
		case PatchDelete:
			working, res = applyDelete(working, inst)
		case PatchRename:
			working, res = applyRename(working, inst)
		default:
			return nil, nil, fmt.Errorf("patcher: unknown op %q", inst.Op)
		}
		results = append(results, res)
	}
	return working, results, nil
}

func applySet(entries []Entry, inst PatchInstruction) ([]Entry, PatchResult) {
	for i, e := range entries {
		if e.Key == inst.Key {
			entries[i].Value = inst.Value
			return entries, PatchResult{Instruction: inst, Applied: true, Note: "updated"}
		}
	}
	entries = append(entries, Entry{Key: inst.Key, Value: inst.Value})
	return entries, PatchResult{Instruction: inst, Applied: true, Note: "inserted"}
}

func applyDelete(entries []Entry, inst PatchInstruction) ([]Entry, PatchResult) {
	for i, e := range entries {
		if e.Key == inst.Key {
			entries = append(entries[:i], entries[i+1:]...)
			return entries, PatchResult{Instruction: inst, Applied: true, Note: "deleted"}
		}
	}
	return entries, PatchResult{Instruction: inst, Applied: false, Note: "key not found"}
}

func applyRename(entries []Entry, inst PatchInstruction) ([]Entry, PatchResult) {
	for i, e := range entries {
		if e.Key == inst.Key {
			entries[i].Key = inst.NewKey
			return entries, PatchResult{Instruction: inst, Applied: true, Note: "renamed"}
		}
	}
	return entries, PatchResult{Instruction: inst, Applied: false, Note: "key not found"}
}
