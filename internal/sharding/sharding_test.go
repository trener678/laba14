package sharding

import (
	"reflect"
	"testing"

	"laba14-health-pipeline/internal/health"
)

func TestAssignSourcesIsDeterministic(t *testing.T) {
	collectors := []string{"collector-b", "collector-a"}
	sources := health.DemoSources()

	first := AssignSources(collectors, sources)
	second := AssignSources(collectors, sources)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("expected deterministic assignment\nfirst=%v\nsecond=%v", first, second)
	}
}

func TestAssignSourcesCoversEverySourceOnce(t *testing.T) {
	assignments := AssignSources([]string{"a", "b", "c"}, health.DemoSources())
	seen := map[string]int{}
	for _, assignment := range assignments {
		for _, source := range assignment.Sources {
			seen[source.ID]++
		}
	}
	if len(seen) != len(health.DemoSources()) {
		t.Fatalf("expected %d sources, got %d", len(health.DemoSources()), len(seen))
	}
	for id, count := range seen {
		if count != 1 {
			t.Fatalf("source %s assigned %d times", id, count)
		}
	}
}
