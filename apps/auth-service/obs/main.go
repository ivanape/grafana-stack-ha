package obs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grafana/pyroscope-go"
	otelchimetric "github.com/riandyrn/otelchi/metric"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	DefaultServiceTags = map[string]string{
		"service": "auth-service",
		"app":     "example",
		"env":     "development",
	}
)

func NewTracer() (trace.Tracer, error) {
	// create otlp exporter, notice that here we are using insecure option
	// because we just want to export the trace locally, also notice that
	// here we don't set any endpoint because by default the otel will load
	// the endpoint from the environment variable `OTEL_EXPORTER_OTLP_ENDPOINT`
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(otlptracehttp.WithInsecure()),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize exporter due: %w", err)
	}
	// initialize tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(DefaultServiceTags["service"]),
		)),
	)
	// set tracer provider and propagator properly, this is to ensure all
	// instrumentation library could run well
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// returns tracer
	return otel.Tracer(DefaultServiceTags["service"]), nil
}

func NewMetricConfig() (otelchimetric.BaseConfig, error) {
	// create context
	ctx := context.Background()

	// create otlp exporter using HTTP protocol. the endpoint will be loaded from
	// OTEL_EXPORTER_OTLP_METRICS_ENDPOINT environment variable
	exporter, err := otlpmetrichttp.New(
		ctx,
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return otelchimetric.BaseConfig{}, err
	}

	// create resource with service name
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(DefaultServiceTags["service"]),
		),
	)
	if err != nil {
		return otelchimetric.BaseConfig{}, err
	}

	// create meter provider with otlp exporter
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(15*time.Second),
			),
		),
	)

	// set global meter provider
	otel.SetMeterProvider(meterProvider)

	// create and return base config for metrics with meter provider
	return otelchimetric.NewBaseConfig(DefaultServiceTags["service"],
		otelchimetric.WithMeterProvider(meterProvider),
	), nil
}

func NewProfiler() (*pyroscope.Profiler, error) {
	config := pyroscope.Config{
		ApplicationName: DefaultServiceTags["service"],
		ServerAddress:   os.Getenv("PYROSCOPE_SERVER_ADDRESS"),
		Logger:          pyroscope.StandardLogger,
		Tags:            DefaultServiceTags,
	}

	return pyroscope.Start(config)
}

func LogErrorWithSpan(logger *logrus.Logger, span trace.Span, context context.Context, msg ...interface{}) {
	spanID := trace.SpanContextFromContext(context).SpanID().String()
	traceID := trace.SpanContextFromContext(context).TraceID().String()

	logger.WithContext(context).
		WithFields(logrus.Fields{
			"spanID":  spanID,
			"traceID": traceID,
		}).
		Error(msg...)
}

func LogInfoWithSpan(logger *logrus.Logger, span trace.Span, context context.Context, msg ...interface{}) {
	spanID := span.SpanContext().SpanID().String()
	traceID := span.SpanContext().TraceID().String()

	logger.WithContext(context).
		WithFields(logrus.Fields{
			"spanID":  spanID,
			"traceID": traceID,
		}).
		Info(msg...)
}
