package env

import (
	"testing"
)

func makeComparatorEntries(pairs ...string) []Entry {
	var entries []Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestComparator_Identical(t *testing.T) {
	c := NewComparator()
	left := makeComparatorEntries("A", "1", "B", "2")
	right := makeComparatorEntries("A", "1", "B", "2")
	res := c.Compare(left, right)
	if len(res.Identical) != 2 {
		t.Errorf("expected 2 identical, got %d", len(res.Identical))
	}
	if len(res.Changed) != 0 || len(res.OnlyInLeft) != 0 || len(res.OnlyInRight) != 0 {
		t.Errorf("unexpected differences: %s", res.Summary())
	}
}

func TestComparator_OnlyInLeft(t *testing.T) {
	c := NewComparator()
	left := makeComparatorEntries("A", "1", "B", "2")
	right := makeComparatorEntries("A", "1")
	res := c.Compare(left, right)
	if len(res.OnlyInLeft) != 1 || res.OnlyInLeft[0].Key != "B" {
		t.Errorf("expected B only in left, got %+v", res.OnlyInLeft)
	}
}

func TestComparator_OnlyInRight(t *testing.T) {
	c := NewComparator()
	left := makeComparatorEntries("A", "1")
	right := makeComparatorEntries("A", "1", "C", "3")
	res := c.Compare(left, right)
	if len(res.OnlyInRight) != 1 || res.OnlyInRight[0].Key != "C" {
		t.Errorf("expected C only in right, got %+v", res.OnlyInRight)
	}
}

func TestComparator_Changed(t *testing.T) {
	c := NewComparator()
	left := makeComparatorEntries("A", "old", "B", "same")
	right := makeComparatorEntries("A", "new", "B", "same")
	res := c.Compare(left, right)
	if len(res.Changed) != 1 {
		t.Fatalf("expected 1 changed, got %d", len(res.Changed))
	}
	if res.Changed[0].Key != "A" || res.Changed[0].OldValue != "old" || res.Changed[0].NewValue != "new" {
		t.Errorf("unexpected change: %+v", res.Changed[0])
	}
}

func TestComparator_CaseInsensitive(t *testing.T) {
	c := NewComparator(WithCaseInsensitive())
	left := makeComparatorEntries("KEY", "Hello")
	right := makeComparatorEntries("KEY", "hello")
	res := c.Compare(left, right)
	if len(res.Identical) != 1 {
		t.Errorf("expected identical with case-insensitive, got changed: %+v", res.Changed)
	}
}

func TestComparator_Summary(t *testing.T) {
	c := NewComparator()
	left := makeComparatorEntries("A", "1", "B", "old")
	right := makeComparatorEntries("B", "new", "C", "3")
	res := c.Compare(left, right)
	summary := res.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	expected := "identical=0 changed=1 only_in_left=1 only_in_right=1"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}

func TestComparator_EmptyInputs(t *testing.T) {
	c := NewComparator()
	res := c.Compare(nil, nil)
	if len(res.Identical)+len(res.Changed)+len(res.OnlyInLeft)+len(res.OnlyInRight) != 0 {
		t.Errorf("expected empty result for nil inputs, got %s", res.Summary())
	}
}
