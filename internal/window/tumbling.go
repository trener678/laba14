package window

import (
	"math"
	"time"

	"laba14-health-pipeline/internal/health"
)

type TumblingAggregator struct {
	size    time.Duration
	buckets map[bucketKey]*state
}

type bucketKey struct {
	metric string
	start  time.Time
}

type state struct {
	count int64
	sum   float64
	min   float64
	max   float64
}

func NewTumblingAggregator(size time.Duration) *TumblingAggregator {
	if size <= 0 {
		size = 10 * time.Second
	}
	return &TumblingAggregator{
		size:    size,
		buckets: map[bucketKey]*state{},
	}
}

func (a *TumblingAggregator) Add(r health.Reading) []health.Aggregate {
	start := r.CollectedAt.Truncate(a.size)
	key := bucketKey{metric: r.Metric, start: start}
	s, ok := a.buckets[key]
	if !ok {
		s = &state{min: math.Inf(1), max: math.Inf(-1)}
		a.buckets[key] = s
	}
	s.count++
	s.sum += r.Value
	s.min = math.Min(s.min, r.Value)
	s.max = math.Max(s.max, r.Value)
	return a.FlushBefore(start)
}

func (a *TumblingAggregator) FlushBefore(cutoff time.Time) []health.Aggregate {
	out := make([]health.Aggregate, 0)
	for key, s := range a.buckets {
		if key.start.Before(cutoff) {
			out = append(out, toAggregate(key, s, a.size))
			delete(a.buckets, key)
		}
	}
	return out
}

func (a *TumblingAggregator) FlushAll() []health.Aggregate {
	out := make([]health.Aggregate, 0, len(a.buckets))
	for key, s := range a.buckets {
		out = append(out, toAggregate(key, s, a.size))
		delete(a.buckets, key)
	}
	return out
}

func toAggregate(key bucketKey, s *state, size time.Duration) health.Aggregate {
	return health.Aggregate{
		Metric:      key.metric,
		WindowStart: key.start,
		WindowEnd:   key.start.Add(size),
		Count:       s.count,
		Min:         s.min,
		Max:         s.max,
		Sum:         s.sum,
		Average:     s.sum / float64(s.count),
	}
}
