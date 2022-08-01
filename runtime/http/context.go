package http

import (
	"context"
	"net/http"

	"github.com/devil-dwj/wms/runtime"
)

type Context struct {
	R         *http.Request
	ReqHeader headerCarrier
	RspHeader headerCarrier
}

func (c *Context) Request() *http.Request {
	return c.R
}

func (c *Context) Kind() runtime.Kind {
	return runtime.KindHTTP
}

func (c *Context) Endpoint() string {
	return c.R.URL.String()
}

func (c *Context) Operation() string {
	return c.R.URL.Path
}

func (c *Context) RequestHeader() runtime.Header {
	return c.ReqHeader
}

func (c *Context) ReplyHeader() runtime.Header {
	return c.RspHeader
}

type headerCarrier http.Header

func (hc headerCarrier) Get(key string) string {
	return http.Header(hc).Get(key)
}

func (hc headerCarrier) Set(key string, value string) {
	http.Header(hc).Set(key, value)
}

func (hc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(hc))
	for k := range http.Header(hc) {
		keys = append(keys, k)
	}
	return keys
}

func HttpContextFromContext(ctx context.Context) (*Context, bool) {
	c, ok := runtime.ContextFromContext(ctx)
	if !ok {
		return nil, false
	}
	hc, ok := c.(*Context)
	if !ok {
		return nil, false
	}
	return hc, true
}
