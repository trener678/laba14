package arrowipc

import (
	"bytes"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow/ipc"

	"laba14-health-pipeline/internal/health"
)

func TestWriteAggregatesProducesArrowStream(t *testing.T) {
	var buf bytes.Buffer
	err := WriteAggregates(&buf, []health.Aggregate{{
		Metric:      "heart_rate",
		WindowStart: time.Unix(100, 0),
		WindowEnd:   time.Unix(110, 0),
		Count:       2,
		Min:         70,
		Max:         80,
		Sum:         150,
		Average:     75,
	}})
	if err != nil {
		t.Fatalf("write arrow: %v", err)
	}

	reader, err := ipc.NewReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("read arrow: %v", err)
	}
	defer reader.Release()

	if !reader.Next() {
		t.Fatal("expected one record batch")
	}
	record := reader.Record()
	if record.NumRows() != 1 || record.NumCols() != 8 {
		t.Fatalf("unexpected record shape: rows=%d cols=%d", record.NumRows(), record.NumCols())
	}
}
