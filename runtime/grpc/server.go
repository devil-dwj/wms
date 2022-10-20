package grpc

import (
	"context"
	"net"

	"github.com/devil-dwj/wms/middleware"
	"google.golang.org/grpc"
)

type ServerOption func(s *Server)

func Address(addr string) ServerOption {
	return func(s *Server) {
		s.addr = addr
	}
}

func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.chain = m
	}
}

type Server struct {
	*grpc.Server
	addr  string
	lis   net.Listener
	chain []middleware.Middleware
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{}
	for _, opt := range opts {
		opt(srv)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ChainServerUnaryInterceptor(srv.chain),
		),
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	return srv
}

func (s *Server) Run(context.Context) error {
	if err := s.listen(); err != nil {
		return err
	}
	return s.Serve(s.lis)
}

func (s *Server) Stop(context.Context) error {

	return nil

}
func (s *Server) listen() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.addr)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	return nil
}
