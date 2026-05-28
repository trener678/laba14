package arrowipc

import (
	"io"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/ipc"
	"github.com/apache/arrow-go/v18/arrow/memory"

	"laba14-health-pipeline/internal/health"
)

func Schema() *arrow.Schema {
	return arrow.NewSchema([]arrow.Field{
		{Name: "metric", Type: arrow.BinaryTypes.String},
		{Name: "window_start_unix", Type: arrow.PrimitiveTypes.Int64},
		{Name: "window_end_unix", Type: arrow.PrimitiveTypes.Int64},
		{Name: "count", Type: arrow.PrimitiveTypes.Int64},
		{Name: "min", Type: arrow.PrimitiveTypes.Float64},
		{Name: "max", Type: arrow.PrimitiveTypes.Float64},
		{Name: "sum", Type: arrow.PrimitiveTypes.Float64},
		{Name: "average", Type: arrow.PrimitiveTypes.Float64},
	}, nil)
}

func WriteAggregates(w io.Writer, aggregates []health.Aggregate) error {
	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, Schema())
	defer builder.Release()

	for _, aggregate := range aggregates {
		builder.Field(0).(*array.StringBuilder).Append(aggregate.Metric)
		builder.Field(1).(*array.Int64Builder).Append(aggregate.WindowStart.Unix())
		builder.Field(2).(*array.Int64Builder).Append(aggregate.WindowEnd.Unix())
		builder.Field(3).(*array.Int64Builder).Append(aggregate.Count)
		builder.Field(4).(*array.Float64Builder).Append(aggregate.Min)
		builder.Field(5).(*array.Float64Builder).Append(aggregate.Max)
		builder.Field(6).(*array.Float64Builder).Append(aggregate.Sum)
		builder.Field(7).(*array.Float64Builder).Append(aggregate.Average)
	}

	record := builder.NewRecord()
	defer record.Release()

	writer := ipc.NewWriter(w, ipc.WithSchema(record.Schema()))
	defer writer.Close()
	return writer.Write(record)
}
