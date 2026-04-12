package env

import (
	"testing"
)

func makePatchEntries() []Entry {
	return []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "5432"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestPatcher_Set_Update(t *testing.T) {
	p := NewPatcher([]PatchInstruction{
		{Op: PatchSet, Key: "PORT", Value: "9999"},
	})
	out, results, err := p.Apply(makePatchEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Note != "updated" {
		t.Errorf("expected note 'updated', got %q", results[0].Note)
	}
	if out[1].Value != "9999" {
		t.Errorf("expected PORT=9999, got %s", out[1].Value)
	}
}

func TestPatcher_Set_Insert(t *testing.T) {
	p := NewPatcher([]PatchInstruction{
		{Op: PatchSet, Key: "NEW_KEY", Value: "hello"},
	})
	out, results, err := p.Apply(makePatchEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Note != "inserted" {
		t.Errorf("expected note 'inserted', got %q", results[0].Note)
	}
	if len(out) != 4 {
		t.Errorf("expected 4 entries, got %d", len(out))
	}
}

func TestPatcher_Delete_Existing(t *testing.T) {
	p := NewPatcher([]PatchInstruction{
		{Op: PatchDelete, Key: "DEBUG"},
	})
	out, results, err := p.Apply(makePatchEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Error("expected Applied=true")
	}
	if len(out) != 2 {
		t.Errorf("expected 2 entries after delete, got %d", len(out))
	}
}

func TestPatcher_Delete_Missing(t *testing.T) {
	p := NewPatcher([]PatchInstruction{
		{Op: PatchDelete, Key: "MISSING"},
	})
	_, results, err := p.Apply(makePatchEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Error("expected Applied=false for missing key")
	}
}

func TestPatcher_Rename(t *testing.T) {
	p := NewPatcher([]PatchInstruction{
		{Op: PatchRename, Key: "HOST", NewKey: "DB_HOST"},
	})
	out, results, err := p.Apply(makePatchEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Applied {
		t.Error("expected Applied=true")
	}
	if out[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", out[0].Key)
	}
}

func TestPatcher_UnknownOp(t *testing.T) {
	p := NewPatcher([]PatchInstruction{
		{Op: PatchOp("noop"), Key: "X"},
	})
	_, _, err := p.Apply(makePatchEntries())
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatcher_DoesNotMutateOriginal(t *testing.T) {
	orig := makePatchEntries()
	p := NewPatcher([]PatchInstruction{
		{Op: PatchSet, Key: "HOST", Value: "changed"},
	})
	_, _, _ = p.Apply(orig)
	if orig[0].Value != "localhost" {
		t.Error("original entries were mutated")
	}
}
