package grpc

import (
	"context"

	"github.com/devil-dwj/wms/middleware"
	"google.golang.org/grpc"
)

func ChainServerUnaryInterceptor(chain []middleware.Middleware) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		h = middleware.Chain(chain...)(h)
		return h(ctx, req)
	}
}

func ChainClientUnaryInterceptor(chain []middleware.Middleware) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return reply, invoker(ctx, method, req, reply, cc, opts...)
		}
		if len(chain) > 0 {
			h = middleware.Chain(chain...)(h)
		}
		_, err := h(ctx, req)
		return err
	}
}
