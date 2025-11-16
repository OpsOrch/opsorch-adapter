package incident

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	coreincident "github.com/opsorch/opsorch-core/incident"
	"github.com/opsorch/opsorch-core/schema"
)

// ProviderName is the registry key for this sample adapter. Replace it with your provider name.
const ProviderName = "example"

// AdapterVersion and RequiresCore express compatibility.
const (
	AdapterVersion = "0.1.0"
	RequiresCore   = ">=0.1.0"
)

var errNotFound = errors.New("incident not found")

// Config captures decrypted configuration from OpsOrch Core.
type Config struct {
	Source          string
	DefaultSeverity string
}

// ExampleProvider is a minimal in-memory incident provider implementation for reference.
type ExampleProvider struct {
	cfg       Config
	mu        sync.Mutex
	nextID    int
	incidents map[string]schema.Incident
	timeline  map[string][]schema.TimelineEntry
}

// New constructs the provider from decrypted config.
func New(cfg map[string]any) (coreincident.Provider, error) {
	parsed := parseConfig(cfg)
	return &ExampleProvider{
		cfg:       parsed,
		nextID:    0,
		incidents: make(map[string]schema.Incident),
		timeline:  make(map[string][]schema.TimelineEntry),
	}, nil
}

func init() {
	_ = coreincident.RegisterProvider(ProviderName, New)
}

// List returns all incidents currently held in memory.
func (p *ExampleProvider) List(ctx context.Context) ([]schema.Incident, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	out := make([]schema.Incident, 0, len(p.incidents))
	for _, inc := range p.incidents {
		out = append(out, inc)
	}
	return out, nil
}

// Get returns a single incident by ID.
func (p *ExampleProvider) Get(ctx context.Context, id string) (schema.Incident, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	inc, ok := p.incidents[id]
	if !ok {
		return schema.Incident{}, errNotFound
	}
	return inc, nil
}

// Create inserts an incident with generated ID and timestamps.
func (p *ExampleProvider) Create(ctx context.Context, in schema.CreateIncidentInput) (schema.Incident, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.nextID++
	now := time.Now()
	id := fmt.Sprintf("ex-%d", p.nextID)
	incident := schema.Incident{
		ID:        id,
		Title:     in.Title,
		Status:    in.Status,
		Severity:  defaultString(in.Severity, p.cfg.DefaultSeverity),
		Service:   in.Service,
		CreatedAt: now,
		UpdatedAt: now,
		Fields:    cloneMap(in.Fields),
		Metadata:  cloneMap(in.Metadata),
	}
	if incident.Metadata == nil {
		incident.Metadata = map[string]any{}
	}
	incident.Metadata["source"] = p.cfg.Source

	p.incidents[id] = incident
	return incident, nil
}

// Update mutates incident fields and bumps UpdatedAt.
func (p *ExampleProvider) Update(ctx context.Context, id string, in schema.UpdateIncidentInput) (schema.Incident, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	inc, ok := p.incidents[id]
	if !ok {
		return schema.Incident{}, errNotFound
	}

	if in.Title != nil {
		inc.Title = *in.Title
	}
	if in.Status != nil {
		inc.Status = *in.Status
	}
	if in.Severity != nil {
		inc.Severity = *in.Severity
	}
	if in.Service != nil {
		inc.Service = *in.Service
	}
	if in.Fields != nil {
		inc.Fields = cloneMap(in.Fields)
	}
	if in.Metadata != nil {
		inc.Metadata = cloneMap(in.Metadata)
	}
	inc.UpdatedAt = time.Now()

	p.incidents[id] = inc
	return inc, nil
}

// Query filters incidents matching the provided query.
func (p *ExampleProvider) Query(ctx context.Context, q schema.IncidentQuery) ([]schema.Incident, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []schema.Incident
	for _, inc := range p.incidents {
		if !matchesQuery(inc, q) {
			continue
		}
		out = append(out, inc)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	if q.Limit > 0 && len(out) > q.Limit {
		out = out[:q.Limit]
	}
	return out, nil
}

// GetTimeline returns the timeline entries for an incident.
func (p *ExampleProvider) GetTimeline(ctx context.Context, id string) ([]schema.TimelineEntry, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	entries := p.timeline[id]
	out := make([]schema.TimelineEntry, len(entries))
	copy(out, entries)
	return out, nil
}

// AppendTimeline appends a timeline entry to an incident.
func (p *ExampleProvider) AppendTimeline(ctx context.Context, id string, entry schema.TimelineAppendInput) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.incidents[id]; !ok {
		return errNotFound
	}

	n := len(p.timeline[id]) + 1
	p.timeline[id] = append(p.timeline[id], schema.TimelineEntry{
		ID:         fmt.Sprintf("%s-t%d", id, n),
		IncidentID: id,
		At:         entry.At,
		Kind:       entry.Kind,
		Body:       entry.Body,
		Actor:      cloneMap(entry.Actor),
		Metadata:   cloneMap(entry.Metadata),
	})
	return nil
}

func parseConfig(cfg map[string]any) Config {
	out := Config{Source: "example", DefaultSeverity: "sev3"}
	if v, ok := cfg["source"].(string); ok && v != "" {
		out.Source = v
	}
	if v, ok := cfg["defaultSeverity"].(string); ok && v != "" {
		out.DefaultSeverity = v
	}
	return out
}

func cloneMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func defaultString(val string, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}

func matchesQuery(inc schema.Incident, q schema.IncidentQuery) bool {
	if q.Query != "" {
		needle := strings.ToLower(q.Query)
		haystack := strings.ToLower(inc.ID + " " + inc.Title)
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	if len(q.Statuses) > 0 && !containsString(q.Statuses, inc.Status) {
		return false
	}
	if len(q.Severities) > 0 && !containsString(q.Severities, inc.Severity) {
		return false
	}
	if svc := strings.TrimSpace(q.Scope.Service); svc != "" && !strings.EqualFold(inc.Service, svc) {
		return false
	}
	if env := strings.TrimSpace(q.Scope.Environment); env != "" {
		if !metadataEqual(inc.Metadata, "environment", env, true) {
			return false
		}
	}
	if team := strings.TrimSpace(q.Scope.Team); team != "" {
		if !metadataEqual(inc.Metadata, "team", team, true) {
			return false
		}
	}
	if len(q.Metadata) > 0 {
		for k, v := range q.Metadata {
			if !metadataEqual(inc.Metadata, k, v, false) {
				return false
			}
		}
	}

	return true
}

func containsString(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func metadataEqual(meta map[string]any, key string, expect any, foldString bool) bool {
	if meta == nil {
		return false
	}
	val, ok := meta[key]
	if !ok {
		return false
	}
	if foldString {
		expectStr, ok1 := expect.(string)
		valStr, ok2 := val.(string)
		if ok1 && ok2 {
			return strings.EqualFold(valStr, expectStr)
		}
	}
	return val == expect
}
