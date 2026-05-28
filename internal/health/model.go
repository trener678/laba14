package health

import (
	"fmt"
	"math"
	"time"
)

type Reading struct {
	PatientID   string    `json:"patient_id"`
	SensorID    string    `json:"sensor_id"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	CollectedAt time.Time `json:"collected_at"`
}

type Aggregate struct {
	Metric      string    `json:"metric"`
	WindowStart time.Time `json:"window_start"`
	WindowEnd   time.Time `json:"window_end"`
	Count       int64     `json:"count"`
	Min         float64   `json:"min"`
	Max         float64   `json:"max"`
	Sum         float64   `json:"sum"`
	Average     float64   `json:"average"`
}

type SensorSource struct {
	ID        string
	PatientID string
	Metric    string
	Unit      string
	BaseValue float64
}

func ValidateReading(r Reading) error {
	if r.PatientID == "" {
		return fmt.Errorf("patient_id is required")
	}
	if r.SensorID == "" {
		return fmt.Errorf("sensor_id is required")
	}
	if r.Metric == "" {
		return fmt.Errorf("metric is required")
	}
	if r.Unit == "" {
		return fmt.Errorf("unit is required")
	}
	if r.CollectedAt.IsZero() {
		return fmt.Errorf("collected_at is required")
	}
	if math.IsNaN(r.Value) || math.IsInf(r.Value, 0) {
		return fmt.Errorf("value must be finite")
	}
	switch r.Metric {
	case "heart_rate":
		if r.Value < 25 || r.Value > 240 {
			return fmt.Errorf("heart_rate out of range")
		}
	case "spo2":
		if r.Value < 50 || r.Value > 100 {
			return fmt.Errorf("spo2 out of range")
		}
	case "temperature":
		if r.Value < 30 || r.Value > 45 {
			return fmt.Errorf("temperature out of range")
		}
	default:
		return fmt.Errorf("unknown metric %q", r.Metric)
	}
	return nil
}

func DemoSources() []SensorSource {
	return []SensorSource{
		{ID: "hr-001", PatientID: "p-1001", Metric: "heart_rate", Unit: "bpm", BaseValue: 72},
		{ID: "hr-002", PatientID: "p-1002", Metric: "heart_rate", Unit: "bpm", BaseValue: 81},
		{ID: "spo2-001", PatientID: "p-1001", Metric: "spo2", Unit: "%", BaseValue: 97},
		{ID: "spo2-002", PatientID: "p-1003", Metric: "spo2", Unit: "%", BaseValue: 95},
		{ID: "temp-001", PatientID: "p-1002", Metric: "temperature", Unit: "C", BaseValue: 36.7},
		{ID: "temp-002", PatientID: "p-1003", Metric: "temperature", Unit: "C", BaseValue: 37.1},
	}
}
