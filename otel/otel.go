package otel

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type OtelEnv string

const (
	OtelEnvProd OtelEnv = "prod"
	OtelEnvDev  OtelEnv = "dev"
)

const (
	OTEL_SERVICE_NAME = "RollDiceService"
)

var (
	// "https://collector:3030"
	otlpmetrichttp_endpoint = os.Getenv("OTLPMETRICHTTP_ENDPOINT")
)

func ParseEnv(env string) OtelEnv {
	switch env {
	case "prod":
		return OtelEnvProd
	case "dev":
		return OtelEnvDev
	default:
		log.Printf("Unknown env \"%s\", using \"dev\"", env)
		return OtelEnvDev
	}
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, env OtelEnv) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	//
	// Set up trace provider.
	//
	var tracerProvider *sdktrace.TracerProvider
	switch env {
	case OtelEnvDev:
		tracerProvider, err = newTraceProvider()
	case OtelEnvProd:
		tracerProvider, err = newOtlpTraceProvider(ctx, otlpmetrichttp_endpoint)
	default:
		log.Fatal("Unknown env")
	}
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	//
	// Set up meter provider.
	//
	var meterProvider *metric.MeterProvider
	switch env {
	case OtelEnvDev:
		meterProvider, err = newMeterProvider()
	case OtelEnvProd:
		meterProvider, err = newOtlpMetricProvider(ctx, otlpmetrichttp_endpoint)
	default:
		log.Fatal("Unknown env")
	}
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newOtlpTraceProvider(ctx context.Context, endpoint string) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}

	traceRes, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(OTEL_SERVICE_NAME),
			// semconv.ServiceVersion(appVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithResource(traceRes),
	)

	return traceProvider, nil
}

func newTraceProvider() (*sdktrace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			sdktrace.WithBatchTimeout(time.Second)),
	)
	return traceProvider, nil
}

func newOtlpMetricProvider(ctx context.Context, endpoint string) (*metric.MeterProvider, error) {
	exp, err := otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}

	meterRes, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(OTEL_SERVICE_NAME),
			// semconv.ServiceVersion(appVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exp, metric.WithInterval(10*time.Second))),
		metric.WithResource(meterRes),
	)

	return meterProvider, nil
}

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(10*time.Second))),
	)
	return meterProvider, nil
}
