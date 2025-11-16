package incident

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/opsorch/opsorch-core/schema"
)

func TestParseConfigDefaults(t *testing.T) {
	cfg := parseConfig(map[string]any{})
	if cfg.Source != "example" {
		t.Fatalf("expected default source, got %q", cfg.Source)
	}
	if cfg.DefaultSeverity != "sev3" {
		t.Fatalf("expected default severity, got %q", cfg.DefaultSeverity)
	}
}

func TestParseConfigOverride(t *testing.T) {
	cfg := parseConfig(map[string]any{"source": "demo", "defaultSeverity": "sev1"})
	if cfg.Source != "demo" {
		t.Fatalf("expected overridden source, got %q", cfg.Source)
	}
	if cfg.DefaultSeverity != "sev1" {
		t.Fatalf("expected overridden severity, got %q", cfg.DefaultSeverity)
	}

	cfg = parseConfig(map[string]any{"source": "", "defaultSeverity": ""})
	if cfg.Source != "example" || cfg.DefaultSeverity != "sev3" {
		t.Fatalf("empty strings should not override defaults: %+v", cfg)
	}
}

func TestDefaultString(t *testing.T) {
	if got := defaultString("value", "fallback"); got != "value" {
		t.Fatalf("expected provided value, got %q", got)
	}
	if got := defaultString("", "fallback"); got != "fallback" {
		t.Fatalf("expected fallback value, got %q", got)
	}
}

func TestCloneMap(t *testing.T) {
	if cloneMap(nil) != nil {
		t.Fatalf("nil map should stay nil")
	}

	original := map[string]any{"k": "v"}
	cloned := cloneMap(original)

	if len(cloned) != len(original) || cloned["k"] != "v" {
		t.Fatalf("cloned map does not match original: %+v", cloned)
	}
	original["k"] = "changed"
	if cloned["k"] != "v" {
		t.Fatalf("clone should not track original mutations, got %q", cloned["k"])
	}
}

func TestExampleProviderCreateAndGet(t *testing.T) {
	ctx := context.Background()
	provAny, err := New(map[string]any{"source": "demo", "defaultSeverity": "sev2"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	prov := provAny.(*ExampleProvider)

	fields := map[string]any{"k": "v"}
	metadata := map[string]any{"extra": "field"}
	inc, err := prov.Create(ctx, schema.CreateIncidentInput{
		Title:    "title",
		Status:   "open",
		Severity: "",
		Service:  "svc-a",
		Fields:   fields,
		Metadata: metadata,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if inc.ID != "ex-1" {
		t.Fatalf("unexpected incident ID: %s", inc.ID)
	}
	if inc.Severity != "sev2" {
		t.Fatalf("expected default severity fallback, got %s", inc.Severity)
	}
	if inc.Metadata["source"] != "demo" {
		t.Fatalf("expected metadata source from config, got %v", inc.Metadata["source"])
	}
	if inc.Metadata["extra"] != "field" {
		t.Fatalf("expected metadata to include input fields, got %v", inc.Metadata)
	}
	if inc.Service != "svc-a" {
		t.Fatalf("expected service to be set, got %q", inc.Service)
	}
	if inc.Fields["k"] != "v" {
		t.Fatalf("expected fields to be copied, got %v", inc.Fields)
	}
	if inc.CreatedAt.IsZero() || inc.UpdatedAt.IsZero() {
		t.Fatalf("timestamps should be populated")
	}
	if !inc.CreatedAt.Equal(inc.UpdatedAt) {
		t.Fatalf("new incidents should have matching timestamps")
	}

	// Mutate input maps to ensure provider stored clones.
	fields["k"] = "mutated"
	metadata["extra"] = "mutated"
	metadata["source"] = "mutated"

	fetched, err := prov.Get(ctx, inc.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if fetched.Fields["k"] != "v" {
		t.Fatalf("stored fields should not reflect caller mutations, got %v", fetched.Fields)
	}
	if fetched.Metadata["extra"] != "field" || fetched.Metadata["source"] != "demo" {
		t.Fatalf("stored metadata should not reflect caller mutations, got %v", fetched.Metadata)
	}

	list, err := prov.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(list) != 1 || list[0].ID != inc.ID {
		t.Fatalf("expected one incident in list, got %+v", list)
	}

	if _, err := prov.Get(ctx, "missing"); !errors.Is(err, errNotFound) {
		t.Fatalf("expected errNotFound for missing incident, got %v", err)
	}
}

func TestExampleProviderUpdate(t *testing.T) {
	ctx := context.Background()
	provAny, err := New(map[string]any{})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	prov := provAny.(*ExampleProvider)

	base, err := prov.Create(ctx, schema.CreateIncidentInput{
		Title:    "old",
		Status:   "open",
		Severity: "sev3",
		Service:  "svc-old",
		Fields:   map[string]any{"k": "v"},
		Metadata: map[string]any{"meta": "data"},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	newTitle := "new title"
	newStatus := "closed"
	newSeverity := "sev1"
	newService := "svc-new"
	newFields := map[string]any{"other": "field"}
	newMetadata := map[string]any{"meta2": "value"}

	updated, err := prov.Update(ctx, base.ID, schema.UpdateIncidentInput{
		Title:    &newTitle,
		Status:   &newStatus,
		Severity: &newSeverity,
		Service:  &newService,
		Fields:   newFields,
		Metadata: newMetadata,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Title != newTitle || updated.Status != newStatus || updated.Severity != newSeverity || updated.Service != newService {
		t.Fatalf("incident fields not updated: %+v", updated)
	}
	if updated.Fields["other"] != "field" {
		t.Fatalf("fields not replaced: %+v", updated.Fields)
	}
	if updated.Metadata["meta2"] != "value" {
		t.Fatalf("metadata not replaced: %+v", updated.Metadata)
	}
	if !updated.UpdatedAt.After(base.UpdatedAt) {
		t.Fatalf("UpdatedAt should advance")
	}

	// Mutations to caller maps should not leak back.
	newFields["other"] = "mutated"
	newMetadata["meta2"] = "mutated"

	fetched, err := prov.Get(ctx, base.ID)
	if err != nil {
		t.Fatalf("Get returned error after update: %v", err)
	}
	if fetched.Fields["other"] != "field" {
		t.Fatalf("stored fields mutated externally: %+v", fetched.Fields)
	}
	if fetched.Metadata["meta2"] != "value" {
		t.Fatalf("stored metadata mutated externally: %+v", fetched.Metadata)
	}

	if _, err := prov.Update(ctx, "missing", schema.UpdateIncidentInput{}); !errors.Is(err, errNotFound) {
		t.Fatalf("expected errNotFound updating missing incident, got %v", err)
	}
}

func TestExampleProviderQuery(t *testing.T) {
	ctx := context.Background()
	provAny, err := New(map[string]any{"source": "demo"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	prov := provAny.(*ExampleProvider)

	inc1, _ := prov.Create(ctx, schema.CreateIncidentInput{ // ex-1
		Title:    "first incident",
		Status:   "open",
		Severity: "sev2",
		Service:  "svc-a",
		Metadata: map[string]any{"team": "platform", "environment": "prod"},
	})
	inc2, _ := prov.Create(ctx, schema.CreateIncidentInput{ // ex-2
		Title:    "second incident",
		Status:   "closed",
		Severity: "sev3",
		Service:  "svc-b",
		Metadata: map[string]any{"team": "sre", "environment": "staging"},
	})
	inc3, _ := prov.Create(ctx, schema.CreateIncidentInput{ // ex-3
		Title:    "another open incident",
		Status:   "open",
		Severity: "sev1",
		Service:  "svc-a",
		Metadata: map[string]any{"team": "platform", "environment": "prod"},
	})

	results, err := prov.Query(ctx, schema.IncidentQuery{Statuses: []string{"open"}, Scope: schema.QueryScope{Service: "svc-a"}, Limit: 1})
	if err != nil {
		t.Fatalf("Query returned error: %v", err)
	}
	if len(results) != 1 || results[0].ID != inc1.ID {
		t.Fatalf("expected limited results starting with first match, got %+v", results)
	}

	results, err = prov.Query(ctx, schema.IncidentQuery{Query: "second"})
	if err != nil {
		t.Fatalf("Query returned error for text search: %v", err)
	}
	if len(results) != 1 || results[0].ID != inc2.ID {
		t.Fatalf("expected to find second incident by query, got %+v", results)
	}

	results, err = prov.Query(ctx, schema.IncidentQuery{Severities: []string{"sev2", "sev1"}, Scope: schema.QueryScope{Environment: "prod", Team: "platform"}})
	if err != nil {
		t.Fatalf("Query returned error for severity/scope filter: %v", err)
	}
	if len(results) != 2 || results[0].ID != inc1.ID || results[1].ID != inc3.ID {
		t.Fatalf("unexpected filtered results: %+v", results)
	}
}

func TestExampleProviderTimeline(t *testing.T) {
	ctx := context.Background()
	provAny, err := New(map[string]any{})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	prov := provAny.(*ExampleProvider)

	inc, err := prov.Create(ctx, schema.CreateIncidentInput{Title: "title", Status: "open"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	at := time.Now()
	actor := map[string]any{"user": "alice"}
	meta := map[string]any{"note": "meta"}
	appendInput := schema.TimelineAppendInput{At: at, Kind: "note", Body: "first", Actor: actor, Metadata: meta}
	if err := prov.AppendTimeline(ctx, inc.ID, appendInput); err != nil {
		t.Fatalf("AppendTimeline returned error: %v", err)
	}

	actor["user"] = "mutated"
	meta["note"] = "mutated"

	entries, err := prov.GetTimeline(ctx, inc.ID)
	if err != nil {
		t.Fatalf("GetTimeline returned error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one timeline entry, got %d", len(entries))
	}
	entry := entries[0]
	if entry.ID != inc.ID+"-t1" {
		t.Fatalf("unexpected entry ID: %s", entry.ID)
	}
	if !entry.At.Equal(at) || entry.Kind != "note" || entry.Body != "first" {
		t.Fatalf("timeline entry not preserved: %+v", entry)
	}
	if entry.Actor["user"] != "alice" {
		t.Fatalf("actor should be cloned, got %v", entry.Actor)
	}
	if entry.Metadata["note"] != "meta" {
		t.Fatalf("metadata should be cloned, got %v", entry.Metadata)
	}

	entries[0].Body = "mutated"
	again, err := prov.GetTimeline(ctx, inc.ID)
	if err != nil {
		t.Fatalf("GetTimeline returned error: %v", err)
	}
	if again[0].Body != "first" {
		t.Fatalf("timeline slice should be copied on read: %+v", again)
	}

	if err := prov.AppendTimeline(ctx, "missing", appendInput); !errors.Is(err, errNotFound) {
		t.Fatalf("expected errNotFound appending timeline to missing incident, got %v", err)
	}

	empty, err := prov.GetTimeline(ctx, "missing")
	if err != nil {
		t.Fatalf("GetTimeline for missing incident returned error: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected empty timeline for missing incident, got %d entries", len(empty))
	}
}
