package app

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const tracingName = "go-cctl"

var Logger *zap.SugaredLogger

func SetupLogging() func() {
	logLevel := zap.NewAtomicLevel()
	logLevel.SetLevel(zap.InfoLevel)
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.Lock(os.Stdout),
		logLevel,
	),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	Logger = logger.Sugar()

	return func() { Logger.Sync() }
}

var ctx context.Context

func SetupTracing() func() {
	_ctx, cancel := context.WithCancel(context.Background())

	ctx = _ctx

	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost("localhost"),
		jaeger.WithAgentPort("6831"),
	))
	if err != nil {
		Logger.Error(err)
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(tracingName),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo")))
	if err != nil {
		panic(err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource))

	otel.SetTracerProvider(tracerProvider)

	return func() {
		_cancel := cancel
		defer _cancel()

		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := tracerProvider.Shutdown(ctx); err != nil {
			Logger.Errorf("error shutting down the tracer provider: %v", err)
		}
	}
}

func Span(spanName string) (context.Context, trace.Span) {
	return otel.Tracer(tracingName).Start(ctx, spanName)
}

func SpanWithContext(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.Tracer(tracingName).Start(ctx, spanName)
}
