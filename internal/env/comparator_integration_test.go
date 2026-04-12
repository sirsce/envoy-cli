package env

import (
	"strings"
	"testing"
)

func TestComparator_WithParser(t *testing.T) {
	leftSrc := "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=development\n"
	rightSrc := "DB_HOST=prod.example.com\nDB_PORT=5432\nAPP_SECRET=abc123\n"

	p := NewParser()

	leftEntries, err := p.Parse(strings.NewReader(leftSrc))
	if err != nil {
		t.Fatalf("parse left: %v", err)
	}
	rightEntries, err := p.Parse(strings.NewReader(rightSrc))
	if err != nil {
		t.Fatalf("parse right: %v", err)
	}

	c := NewComparator()
	res := c.Compare(leftEntries, rightEntries)

	if len(res.Changed) != 1 || res.Changed[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST changed, got %+v", res.Changed)
	}
	if len(res.Identical) != 1 || res.Identical[0].Key != "DB_PORT" {
		t.Errorf("expected DB_PORT identical, got %+v", res.Identical)
	}
	if len(res.OnlyInLeft) != 1 || res.OnlyInLeft[0].Key != "APP_ENV" {
		t.Errorf("expected APP_ENV only in left, got %+v", res.OnlyInLeft)
	}
	if len(res.OnlyInRight) != 1 || res.OnlyInRight[0].Key != "APP_SECRET" {
		t.Errorf("expected APP_SECRET only in right, got %+v", res.OnlyInRight)
	}
}

func TestComparator_WithFilterAndCompare(t *testing.T) {
	src := "DB_HOST=localhost\nDB_PORT=5432\nAPP_ENV=staging\nAPP_DEBUG=true\n"
	p := NewParser()
	allEntries, err := p.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	f := NewFilter()
	dbEntries := f.ByPrefix(allEntries, "DB_")
	appEntries := f.ByPrefix(allEntries, "APP_")

	c := NewComparator()
	res := c.Compare(dbEntries, appEntries)

	// DB_ and APP_ keys are entirely disjoint
	if len(res.OnlyInLeft) != 2 {
		t.Errorf("expected 2 only in left (DB keys), got %d", len(res.OnlyInLeft))
	}
	if len(res.OnlyInRight) != 2 {
		t.Errorf("expected 2 only in right (APP keys), got %d", len(res.OnlyInRight))
	}
	if len(res.Identical)+len(res.Changed) != 0 {
		t.Errorf("expected no overlap, got summary: %s", res.Summary())
	}
}
