package o11y

import "go.opentelemetry.io/otel"

// Tracer is the shared OpenTelemetry tracer used for all gravity spans.
// When no TracerProvider is configured (OTel disabled) it produces no-op spans.
var Tracer = otel.Tracer("beryju.io/gravity")
