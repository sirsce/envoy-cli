package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestAuditor_Integration_PushDiff(t *testing.T) {
	base := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "SECRET", Value: "old-secret"},
	}
	updated := []Entry{
		{Key: "DB_HOST", Value: "prod.db.example.com"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "API_KEY", Value: "new-key"},
	}

	auditor := newTestAuditor()
	events := auditor.Audit(base, updated)
	log := NewAuditLog(events)

	if !log.HasEvents() {
		t.Fatal("expected audit events")
	}

	added := log.FilterByAction(AuditActionAdd)
	removed := log.FilterByAction(AuditActionRemove)
	changed := log.FilterByAction(AuditActionChange)

	if len(added) != 1 || added[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY added, got %+v", added)
	}
	if len(removed) != 1 || removed[0].Key != "SECRET" {
		t.Errorf("expected SECRET removed, got %+v", removed)
	}
	if len(changed) != 1 || changed[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST changed, got %+v", changed)
	}
}

func TestAuditor_Integration_TextReport(t *testing.T) {
	base := []Entry{{Key: "FOO", Value: "1"}}
	updated := []Entry{{Key: "FOO", Value: "2"}, {Key: "BAR", Value: "3"}}

	auditor := newTestAuditor()
	log := NewAuditLog(auditor.Audit(base, updated))

	var buf bytes.Buffer
	if err := log.WriteText(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "CHANGE FOO") {
		t.Errorf("expected CHANGE FOO in output: %q", out)
	}
	if !strings.Contains(out, "ADD BAR") {
		t.Errorf("expected ADD BAR in output: %q", out)
	}
}

func TestAuditor_Integration_JSONReport(t *testing.T) {
	base := []Entry{{Key: "X", Value: "1"}}
	updated := []Entry{}

	auditor := newTestAuditor()
	log := NewAuditLog(auditor.Audit(base, updated))

	var buf bytes.Buffer
	if err := log.WriteJSON(&buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "remove") {
		t.Errorf("expected remove action in JSON: %q", buf.String())
	}
}
