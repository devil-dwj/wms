package tracing

import (
	"context"
	"fmt"

	"github.com/devil-dwj/wms/middleware"
	"github.com/devil-dwj/wms/runtime"
	"github.com/devil-dwj/wms/runtime/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "github.com/devil-dwj/wms/middleware/tracing"
)

type Option func(*options)

type options struct {
	serviceName string
	tp          trace.TracerProvider
}

func WithServiceName(name string) Option {
	return func(o *options) {
		o.serviceName = name
	}
}

func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(o *options) {
		o.tp = tp
	}
}

func Tracing(opts ...Option) middleware.Middleware {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	if o.serviceName == "" {
		o.serviceName = "HTTP SERVICE"
	}
	if o.tp == nil {
		o.tp = otel.GetTracerProvider()
	}
	tracer := o.tp.Tracer(tracerName, trace.WithInstrumentationVersion("semver: 0.0.0"))
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			rtCtx, ok := runtime.ContextFromContext(ctx)
			if !ok {
				return h(ctx, req)
			}
			httpCtx, ok := rtCtx.(*http.Context)
			if !ok {
				return h(ctx, req)
			}
			traceOpts := []trace.SpanStartOption{
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", httpCtx.Request())...),
				trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(httpCtx.Request())...),
				trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(o.serviceName, httpCtx.Request().URL.Path, httpCtx.Request())...),
			}
			spanName := httpCtx.Request().URL.Path
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s router not found", httpCtx.Request().Method)
			}
			var span trace.Span
			ctx, span = tracer.Start(ctx, spanName, traceOpts...)
			defer span.End()

			reply, err := h(ctx, req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			return reply, err
		}
	}
}
