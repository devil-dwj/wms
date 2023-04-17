package rabbitmq

import (
	"context"
	"time"

	"github.com/streadway/amqp"
)

type Option func(*options)

type options struct {
	heartbeat time.Duration
	tracer    bool
}

func WithHeartBeat(t time.Duration) Option {
	return func(o *options) {
		o.heartbeat = t
	}
}

func WithTracer(b bool) Option {
	return func(o *options) {
		o.tracer = b
	}
}

type Producer interface {
	Publish(ctx context.Context, b []byte, des string) error
}

type producer struct {
	url          string
	exchangeName string
	routingKey   string
	opt          options
	ch           *amqp.Channel
	tracing      Tracing
}

func NewProducer(url string, exchangeName string, routingKey string, opts ...Option) Producer {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	if o.heartbeat == 0 {
		o.heartbeat = 10 * time.Second
	}
	p := &producer{url: url, exchangeName: exchangeName, routingKey: routingKey, opt: o}
	if o.tracer {
		p.tracing = NewTracing(
			WithType("producer"),
		)
	}

	c, err := amqp.DialConfig(url, amqp.Config{
		Heartbeat: p.opt.heartbeat,
		Locale:    "en_US",
	})
	if err != nil {
		panic(err)
	}
	ch, err := c.Channel()
	if err != nil {
		panic(err)
	}
	p.ch = ch

	return p
}

func (p *producer) Publish(ctx context.Context, b []byte, des string) error {
	var c context.Context
	if p.tracing != nil {
		c, _ = p.tracing.Before(ctx, des)
	}
	err := p.ch.Publish(
		p.exchangeName,
		p.routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            b,
			DeliveryMode:    0,
			Priority:        0,
		},
	)
	if p.tracing != nil {
		p.tracing.After(c, err)
	}
	return err
}
