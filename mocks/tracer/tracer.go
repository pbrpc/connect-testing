//revive:disable:package-comments
package tracer

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// Mock wraps OTEL's in-memory exporter for testing span events.
type Mock struct {
	exporter *tracetest.InMemoryExporter
	provider *sdktrace.TracerProvider
	span     sdktrace.ReadWriteSpan
}

// New creates a new mock tracer and returns it along with a context containing an active span.
func New(t *testing.T) (*Mock, context.Context) {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	tracer := tp.Tracer("test")
	ctx, span := tracer.Start(t.Context(), "test-span")
	return &Mock{
		exporter: exporter,
		provider: tp,
		span:     span.(sdktrace.ReadWriteSpan),
	}, ctx
}

// EndSpanAndGetEvents ends the current span and returns all recorded events.
func (m *Mock) EndSpanAndGetEvents(t *testing.T) []sdktrace.Event {
	t.Helper()
	m.span.End()
	spans := m.exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected spans")
	}
	return spans[0].Events
}

// EndSpan ends the current span without returning events.
func (m *Mock) EndSpan() {
	m.span.End()
}

// GetSpans returns all recorded spans.
func (m *Mock) GetSpans() tracetest.SpanStubs {
	return m.exporter.GetSpans()
}

// Tracer returns a tracer from the mock's provider.
func (m *Mock) Tracer(name string) trace.Tracer {
	return m.provider.Tracer(name)
}

// Shutdown shuts down the trace provider.
func (m *Mock) Shutdown(t *testing.T) {
	t.Helper()
	_ = m.provider.Shutdown(t.Context())
}
