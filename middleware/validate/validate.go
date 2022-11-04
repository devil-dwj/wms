package validate

import (
	"context"

	"github.com/devil-dwj/wms/middleware"
)

func Validator() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if v, ok := req.(interface {
				ValidateAll() error
			}); ok {
				if err := v.ValidateAll(); err != nil {
					return nil, err
				}
			}
			return h(ctx, req)
		}
	}
}
