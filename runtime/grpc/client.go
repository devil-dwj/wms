package grpc

import (
	"context"

	"github.com/devil-dwj/wms/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientOption func(o *clientOptions)

type clientOptions struct {
	endpoint string
	chain    []middleware.Middleware
}

func WithEndPoint(endpoint string) ClientOption {
	return func(o *clientOptions) {
		o.endpoint = endpoint
	}
}

func WithMiddleware(m ...middleware.Middleware) ClientOption {
	return func(o *clientOptions) {
		o.chain = m
	}
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, opts...)
}

func dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	o := clientOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(
			ChainClientUnaryInterceptor(o.chain),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return grpc.DialContext(ctx, o.endpoint, grpcOpts...)
}
