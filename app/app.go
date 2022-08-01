package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/devil-dwj/wms/runtime"
	"golang.org/x/sync/errgroup"
)

type Option func(o *options)

type options struct {
	ctx     context.Context
	sigs    []os.Signal
	servers []runtime.Server
}

func Server(srv ...runtime.Server) Option {
	return func(o *options) { o.servers = srv }
}

type App struct {
	opts options
}

func New(opts ...Option) *App {
	o := options{
		ctx: context.Background(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	return &App{
		opts: o,
	}
}

func (a *App) Run() error {
	eg, ctx := errgroup.WithContext(a.opts.ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx, cancel := context.WithTimeout(a.opts.ctx, time.Minute)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Run(a.opts.ctx)
		})
	}
	wg.Wait()

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return nil
	}

	return nil
}

func (a *App) Stop() error {
	return nil
}
