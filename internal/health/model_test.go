package health

import (
	"testing"
	"time"
)

func TestValidateReadingRejectsOutOfRange(t *testing.T) {
	r := Reading{
		PatientID:   "p1",
		SensorID:    "s1",
		Metric:      "spo2",
		Value:       120,
		Unit:        "%",
		CollectedAt: time.Now(),
	}
	if err := ValidateReading(r); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestSimulatorProducesValidDemoReadings(t *testing.T) {
	sim := NewSimulator(42)
	for _, src := range DemoSources() {
		r := sim.Reading(src, time.Unix(100, 0))
		if err := ValidateReading(r); err != nil {
			t.Fatalf("reading for %s is invalid: %v", src.ID, err)
		}
	}
}
