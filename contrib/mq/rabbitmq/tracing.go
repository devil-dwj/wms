package rabbitmq

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultTracerName = "github.com/devil-dwj/wms-bee/mq/rabbitmq"
	systemkey         = "rabbitmq"
)

type TracingOption func(*TracingOptions)

type TracingOptions struct {
	tp    trace.TracerProvider
	attrs []attribute.KeyValue
	mtype string
}

func WithTracerProvicer(tp trace.TracerProvider) TracingOption {
	return func(to *TracingOptions) {
		if tp != nil {
			to.tp = tp
		}
	}
}

func WithAttributes(attrs ...attribute.KeyValue) TracingOption {
	return func(to *TracingOptions) {
		to.attrs = append(to.attrs, attrs...)
	}
}

func WithType(t string) TracingOption {
	return func(to *TracingOptions) {
		to.mtype = t
	}
}

type Tracing interface {
	Before(ctx context.Context, destination string) (context.Context, error)
	After(ctx context.Context, err error)
}

type tracing struct {
	opt    TracingOptions
	tracer trace.Tracer
}

func NewTracing(opts ...TracingOption) Tracing {
	o := TracingOptions{
		tp: otel.GetTracerProvider(),
		attrs: []attribute.KeyValue{
			semconv.MessagingSystemKey.String(systemkey),
		},
	}
	for _, opt := range opts {
		opt(&o)
	}
	tracer := o.tp.Tracer(
		defaultTracerName,
	)
	return &tracing{tracer: tracer, opt: o}
}

func (c *tracing) Before(ctx context.Context, destination string) (context.Context, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}
	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(c.opt.attrs...),
		trace.WithAttributes(semconv.MessagingDestinationKey.String(destination)),
	}
	ctx, _ = c.tracer.Start(ctx, c.opt.mtype, opts...)
	return ctx, nil
}

func (c *tracing) After(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
