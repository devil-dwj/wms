package validate

import (
	"context"

	"github.com/devil-dwj/wms/middleware"
)

type validator interface {
	Validate() error
}

func Validator() middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					return nil, err
				}
			}
			return h(ctx, req)
		}
	}
}
