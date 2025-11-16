package main

import (
	"context"
	"testing"

	"github.com/opsorch/opsorch-core/schema"
)

type stubProvider struct{}

func (stubProvider) Query(ctx context.Context, query schema.IncidentQuery) ([]schema.Incident, error) {
	return nil, nil
}
func (stubProvider) List(ctx context.Context) ([]schema.Incident, error) { return nil, nil }
func (stubProvider) Get(ctx context.Context, id string) (schema.Incident, error) {
	return schema.Incident{}, nil
}
func (stubProvider) Create(ctx context.Context, in schema.CreateIncidentInput) (schema.Incident, error) {
	return schema.Incident{}, nil
}
func (stubProvider) Update(ctx context.Context, id string, in schema.UpdateIncidentInput) (schema.Incident, error) {
	return schema.Incident{}, nil
}
func (stubProvider) GetTimeline(ctx context.Context, id string) ([]schema.TimelineEntry, error) {
	return nil, nil
}
func (stubProvider) AppendTimeline(ctx context.Context, id string, entry schema.TimelineAppendInput) error {
	return nil
}

func TestEnsureProviderReturnsExisting(t *testing.T) {
	t.Cleanup(func() { provider = nil })
	existing := stubProvider{}
	provider = existing

	got, err := ensureProvider(map[string]any{"source": "ignored"})
	if err != nil {
		t.Fatalf("ensureProvider returned error: %v", err)
	}
	if got != existing {
		t.Fatalf("expected existing provider to be reused")
	}
}

func TestEnsureProviderCachesNewInstance(t *testing.T) {
	t.Cleanup(func() { provider = nil })

	first, err := ensureProvider(map[string]any{"source": "demo"})
	if err != nil {
		t.Fatalf("ensureProvider returned error: %v", err)
	}
	second, err := ensureProvider(map[string]any{"source": "other"})
	if err != nil {
		t.Fatalf("ensureProvider returned error on second call: %v", err)
	}
	if first != second {
		t.Fatalf("expected provider instance to be cached between calls")
	}

	if first == nil {
		t.Fatalf("expected provider instance to be non-nil")
	}
}
