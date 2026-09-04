package instance

import (
	"context"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

var otelTracerProvider *sdktrace.TracerProvider

func (i *Instance) startOTel() {
	if !extconfig.Get().Observability.OTel.Enabled || extconfig.Get().CI {
		return
	}
	exp, err := otlptracegrpc.New(i.rootContext)
	if err != nil {
		i.log.Warn("failed to init otel exporter", zap.Error(err))
		return
	}
	rate := 0.5
	if extconfig.Get().Debug {
		rate = 1
	}
	otelTracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(rate))),
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", "gravity"),
			attribute.String("service.version", extconfig.FullVersion()),
			attribute.String("gravity.instance", extconfig.Get().Instance.Identifier),
		)),
	)
	otel.SetTracerProvider(otelTracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func (i *Instance) stopOTel() {
	if otelTracerProvider == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := otelTracerProvider.Shutdown(ctx); err != nil {
		i.log.Warn("failed to shutdown otel", zap.Error(err))
	}
}
