package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
/*
* jaegerEntryPoint "http://10.50.32.48:4318/v1/traces"
* jaegerEntryPoint "http://127.0.0.1:14268/v1/trace"
* jaegerEntryPoint "http://127.0.0.1:14268/api/traces"
* environment: prod
 */
func StartOpenTelemetry(serviceName, jaegerEntryPoint string) (*trace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEntryPoint)))
	if err != nil {
		return nil, err
	}

	// forward trace id from client
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// initial new tracer provider
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			// attribute.String("environment", environment),
			// attribute.Int64("ID", id),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	// set global info
	// setGlobalTracer(tp.Tracer(serviceName))

	return tp, nil
}
