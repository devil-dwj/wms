package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/extra/redisotel/v8"
	red "github.com/go-redis/redis/v8"
)

type Option func(*options)

type options struct {
	db       int
	password string
	trace    bool
}

func WithDb(i int) Option {
	return func(o *options) {
		o.db = i
	}
}

func WithPassword(p string) Option {
	return func(o *options) {
		o.password = p
	}
}

func WithTrace(b bool) Option {
	return func(o *options) {
		o.trace = b
	}
}

type Basic interface {
	Exists(ctx context.Context, keys ...string) (bool, error)
	ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)

	Set(ctx context.Context, key string, value string, t time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) (bool, error)

	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	HGet(ctx context.Context, key, field string) (string, error)
	HDel(ctx context.Context, key, field string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error)
}

type basicRedis struct {
	*red.Client
	addr string
	opt  options
}

func NewBasicRedis(addr string, opts ...Option) Basic {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}

	c := red.NewClient(&red.Options{
		Addr:     addr,
		DB:       o.db,
		Password: o.password,
	})
	if o.trace {
		c.AddHook(redisotel.NewTracingHook())
	}
	return &basicRedis{Client: c, addr: addr, opt: o}
}

func (r *basicRedis) Exists(ctx context.Context, keys ...string) (bool, error) {
	i, err := r.Client.Exists(ctx, keys...).Result()
	if err != nil {
		return false, err
	}

	return i == 1, nil
}

func (r *basicRedis) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	return r.Client.ExpireAt(ctx, key, tm).Result()
}

func (r *basicRedis) Set(ctx context.Context, key string, value string, t time.Duration) (bool, error) {
	reply, err := r.Client.Set(ctx, key, value, t).Result()

	return reply == "OK", err
}

func (r *basicRedis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err == red.Nil {
		return "", nil
	}
	return value, err
}

func (r *basicRedis) Del(ctx context.Context, key ...string) (bool, error) {
	i, err := r.Client.Del(ctx, key...).Result()
	if err != nil {
		return false, err
	}

	return i == 1, nil
}

func (r *basicRedis) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.Client.HSet(ctx, key, values...).Result()

	return v, err
}

func (r *basicRedis) HGet(ctx context.Context, key, field string) (string, error) {
	v, err := r.Client.HGet(ctx, key, field).Result()
	if err == red.Nil {
		return v, nil
	}
	return v, err
}

func (r *basicRedis) HDel(ctx context.Context, key, field string) error {
	err := r.Client.HDel(ctx, key, field).Err()
	if err == red.Nil {
		return nil
	}
	return err
}

func (r *basicRedis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	v, err := r.Client.HGetAll(ctx, key).Result()
	if err == red.Nil {
		return map[string]string{}, err
	}
	return v, err
}

func (r *basicRedis) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	v, err := r.Client.HIncrBy(ctx, key, field, incr).Result()
	return v, err
}
