package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Config interface {
	Load() error
	Scan(v interface{}) error
}

type Option func(*options)

type options struct {
	source string
}

func WithSource(s string) Option {
	return func(o *options) {
		o.source = s
	}
}

type config struct {
	opt options
	b   []byte
}

func New(opts ...Option) Config {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &config{
		opt: o,
		b:   make([]byte, 0),
	}
}

func (c *config) Load() error {
	b, err := ioutil.ReadFile(c.opt.source)
	if err != nil {
		return err
	}
	c.b = append(c.b, b...)
	return nil
}

func (c *config) Scan(v interface{}) error {
	d := json.NewDecoder(bytes.NewReader(c.b))
	d.UseNumber()
	return d.Decode(v)
}
