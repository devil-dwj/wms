package http

import (
	"context"
	"net/http"

	"github.com/devil-dwj/wms/log"
	"github.com/devil-dwj/wms/middleware"
	"github.com/devil-dwj/wms/runtime"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type ServerOption func(*Server)

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func Middleware(c ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.chain = c
	}
}

func Static(path string) ServerOption {
	return func(s *Server) {
		s.static = path
	}
}

type ServiceRegistrar interface {
	RegisterService(desc *RouterDesc, impl interface{})
}

type methodHandler func(
	srv interface{},
	ctx context.Context,
	dec func(interface{}) error,
	interceptor ServerInterceptor,
) (interface{}, error)

type RouterDesc struct {
	ServiceName string
	Methods     []MethodDesc
}

type MethodDesc struct {
	Name    string
	Method  string
	Path    string
	Handler methodHandler
}

type routerInfo struct {
	serviceName string
	serveImpl   interface{}
	methods     map[string]*MethodDesc
}

type Server struct {
	*gin.Engine
	addr     string
	unaryInt ServerInterceptor
	chain    []middleware.Middleware
	routers  map[string]*routerInfo
	static   string
}

func NewServer(addr string, opts ...ServerOption) *Server {
	svr := &Server{
		Engine:  gin.New(),
		addr:    addr,
		routers: make(map[string]*routerInfo),
	}

	for _, o := range opts {
		o(svr)
	}
	if svr.static != "" {
		svr.Engine.StaticFS(svr.static, http.Dir(svr.static))
	}

	svr.unaryInt = ChainIterceptor(svr.chain)

	return svr
}

func (s *Server) RegisterService(desc *RouterDesc, impl interface{}) {
	info := &routerInfo{
		serviceName: desc.ServiceName,
		serveImpl:   impl,
		methods:     make(map[string]*MethodDesc),
	}
	for i := range desc.Methods {
		d := &desc.Methods[i]
		if d.Path[:1] != "/" {
			d.Path = "/" + d.Path
		}
		info.methods[d.Path] = d
		s.Handle(d.Method, d.Path, s.handler)
	}
	s.routers[desc.ServiceName] = info
}

func (s *Server) Run(ctx context.Context) error {
	log.Infof("http server start running: %s", s.addr)
	return s.Engine.Run(s.addr)
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) handler(ctx *gin.Context) {
	var df func(interface{}) error = nil
	path := ctx.Request.URL.Path
	for _, info := range s.routers {
		if md, ok := info.methods[path]; ok {
			if md.Method == http.MethodPost {
				df = s.bindBody(ctx)
			} else if md.Method == http.MethodGet {
				df = s.bindQuery(ctx)
			}
			c := &Context{
				R:         ctx.Request,
				ReqHeader: headerCarrier(ctx.Request.Header),
				RspHeader: headerCarrier(ctx.Writer.Header()),
			}
			passCtx := runtime.NewServerContext(ctx.Request.Context(), c)
			reply, err := md.Handler(info.serveImpl, passCtx, df, s.unaryInt)
			if err != nil {
				fail(ctx, err)
			} else {
				success(ctx, reply)
			}
		}
	}
}

func (s *Server) bindBody(ctx *gin.Context) func(interface{}) error {
	return func(i interface{}) error {
		return ctx.ShouldBind(i)
	}
}

func (s *Server) bindQuery(ctx *gin.Context) func(interface{}) error {
	return func(i interface{}) error {
		values := ctx.Request.URL.Query()
		return binding.MapFormWithTag(i, values, "json")
	}
}

func fail(c *gin.Context, err error) {
	var status = http.StatusBadRequest
	var code = 1
	if e, ok := err.(interface {
		Status() int
		Code() int
	}); ok {
		status = e.Status()
		code = e.Code()
	}
	c.Error(err)
	c.JSON(
		status,
		gin.H{
			"code": code,
			"msg":  err.Error(),
			"data": "",
		})
}

func success(c *gin.Context, data interface{}) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"code": 0,
			"msg":  "",
			"data": data,
		})
}

type requestKey struct{}

func NewRequestContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestKey{}, req)
}

func RequestFromContext(ctx context.Context) (req *http.Request, ok bool) {
	req, ok = ctx.Value(requestKey{}).(*http.Request)
	return
}
