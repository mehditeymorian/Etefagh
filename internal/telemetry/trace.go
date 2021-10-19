package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"log"
)

func New(cfg Config) trace.Tracer {
	var exporter sdktrace.SpanExporter

	var err error
	// trace server is listening
	if cfg.Enabled {
		//exporter, err = jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(cfg.Agent.Host), jaeger.WithAgentPort(cfg.Agent.Port)))

		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Url)))
	} else {
		// export traces to standard output
		exporter, err = stdout.New(stdout.WithPrettyPrint())
	}

	if err != nil {
		log.Fatalf("failed to initialize export pipeline: %w", err)
	}

	// use default resource and update server namespaceKey and nameKey
	res, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaless(semconv.ServiceNamespaceKey.String("mehditeymorian"), semconv.ServiceNameKey.String("etefagh")))
	if err != nil {
		panic(err)
	}

	// create a span processor to send the completed spans to the exporter
	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	// create a trace provider with bsp and customized resource
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(bsp), sdktrace.WithResource(res))

	otel.SetTracerProvider(tp)

	tracer := otel.Tracer("timurid.ir/etefagh")

	return tracer
}
