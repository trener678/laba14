package sharding

import (
	"hash/fnv"
	"sort"

	"laba14-health-pipeline/internal/health"
)

type Assignment struct {
	CollectorID string
	Sources     []health.SensorSource
}

func AssignSources(collectors []string, sources []health.SensorSource) map[string]Assignment {
	result := make(map[string]Assignment, len(collectors))
	if len(collectors) == 0 {
		return result
	}

	ids := append([]string(nil), collectors...)
	sort.Strings(ids)
	for _, id := range ids {
		result[id] = Assignment{CollectorID: id}
	}

	for _, source := range sources {
		owner := ids[int(hash(source.ID))%len(ids)]
		assignment := result[owner]
		assignment.Sources = append(assignment.Sources, source)
		result[owner] = assignment
	}
	return result
}

func hash(value string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(value))
	return h.Sum32()
}
