package health

import (
	"math"
	"math/rand"
	"time"
)

type Simulator struct {
	rng *rand.Rand
}

func NewSimulator(seed int64) *Simulator {
	return &Simulator{rng: rand.New(rand.NewSource(seed))}
}

func (s *Simulator) Reading(source SensorSource, at time.Time) Reading {
	noise := s.rng.NormFloat64()
	value := source.BaseValue
	switch source.Metric {
	case "heart_rate":
		value += noise * 5
	case "spo2":
		value += noise * 1.2
	case "temperature":
		value += noise * 0.25
	}
	return Reading{
		PatientID:   source.PatientID,
		SensorID:    source.ID,
		Metric:      source.Metric,
		Value:       round(value, 2),
		Unit:        source.Unit,
		CollectedAt: at.UTC(),
	}
}

func round(v float64, places int) float64 {
	pow := math.Pow10(places)
	return math.Round(v*pow) / pow
}
