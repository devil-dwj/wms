package runtime

import "context"

type Server interface {
	Run(context.Context) error
	Stop(context.Context) error
}

type Header interface {
	Get(key string) string
	Set(key string, value string)
	Keys() []string
}

type Context interface {
	Kind() Kind
	Endpoint() string
	Operation() string
	RequestHeader() Header
	ReplyHeader() Header
}

type Kind string

func (k Kind) String() string { return string(k) }

const (
	KindGRPC Kind = "grpc"
	KindHTTP Kind = "http"
)

type (
	serverContextKey struct{}
)

func NewServerContext(ctx context.Context, c Context) context.Context {
	return context.WithValue(ctx, serverContextKey{}, c)
}

func ContextFromContext(ctx context.Context) (c Context, ok bool) {
	c, ok = ctx.Value(serverContextKey{}).(Context)
	return
}
