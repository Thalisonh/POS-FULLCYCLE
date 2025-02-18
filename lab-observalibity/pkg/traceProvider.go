package pkg

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

func InitProvider(serviceName, colectorURL string, exporter *zipkin.Exporter) (func(context.Context) error, error) {
	ctx := context.Background()

	res, err := resource.New(ctx, resource.WithAttributes(
		semconv.ServiceName(serviceName),
	),
	)
	if err != nil {
		return nil, errors.Join(err, errors.New("fail to start new resource"))
	}

	conn, err := grpc.NewClient(colectorURL, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Join(err, errors.New("fail to start new grpc client"))
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, errors.Join(err, errors.New("fail to start traceExporter"))
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}
