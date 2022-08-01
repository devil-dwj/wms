package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/devil-dwj/wms/log"
	"github.com/devil-dwj/wms/middleware"
	"github.com/devil-dwj/wms/runtime"
)

func Logging(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		rctx, ok := runtime.ContextFromContext(ctx)
		if !ok {
			log.Error("not find runtime context")
		}
		startTime := time.Now()
		reply, err := handler(ctx, req)
		log.LogWithOptions(
			log.WithLevel(log.LevelInfo),
			log.WithSkip(3),
			log.WithKeyVals(
				"path", rctx.Operation(),
				"req", extractArgs(req),
				"reply", extractArgs(reply),
				"cost", time.Since(startTime).Milliseconds(),
			),
		)
		return reply, err
	}
}

func extractArgs(req interface{}) string {
	if s, ok := req.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", req)
}
