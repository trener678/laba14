package window

import (
	"testing"
	"time"

	"laba14-health-pipeline/internal/health"
)

func TestTumblingAggregatorFlushesClosedWindow(t *testing.T) {
	agg := NewTumblingAggregator(10 * time.Second)
	base := time.Unix(100, 0).UTC()

	agg.Add(reading("heart_rate", 70, base.Add(1*time.Second)))
	flushed := agg.Add(reading("heart_rate", 80, base.Add(12*time.Second)))

	if len(flushed) != 1 {
		t.Fatalf("expected one flushed aggregate, got %d", len(flushed))
	}
	got := flushed[0]
	if got.Count != 1 || got.Average != 70 || got.WindowStart != base {
		t.Fatalf("unexpected aggregate: %+v", got)
	}
}

func TestTumblingAggregatorComputesStats(t *testing.T) {
	agg := NewTumblingAggregator(10 * time.Second)
	base := time.Unix(100, 0).UTC()

	agg.Add(reading("spo2", 96, base.Add(time.Second)))
	agg.Add(reading("spo2", 98, base.Add(2*time.Second)))
	out := agg.FlushAll()

	if len(out) != 1 {
		t.Fatalf("expected one aggregate, got %d", len(out))
	}
	if out[0].Count != 2 || out[0].Min != 96 || out[0].Max != 98 || out[0].Average != 97 {
		t.Fatalf("bad aggregate: %+v", out[0])
	}
}

func reading(metric string, value float64, at time.Time) health.Reading {
	return health.Reading{
		PatientID:   "p1",
		SensorID:    "s1",
		Metric:      metric,
		Value:       value,
		Unit:        "u",
		CollectedAt: at,
	}
}
