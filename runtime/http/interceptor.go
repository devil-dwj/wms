package http

import (
	"context"

	"github.com/devil-dwj/wms/middleware"
)

type InterceptorHandler func(ctx context.Context, req interface{}) (interface{}, error)

type ServerInterceptor func(ctx context.Context, req interface{}, handler InterceptorHandler) (resp interface{}, err error)

func ChainIterceptor(chain []middleware.Middleware) ServerInterceptor {
	return func(ctx context.Context, req interface{}, handler InterceptorHandler) (resp interface{}, err error) {
		h := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}
		h = middleware.Chain(chain...)(h)
		return h(ctx, req)
	}
}
