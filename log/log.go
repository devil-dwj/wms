package log

import "log"

var DefaultLogger = NewStdLogger(log.Writer())

type Option func(*Options)

type Options struct {
	Level   Level
	Skip    int
	Keyvals []interface{}
}

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithSkip(i int) Option {
	return func(o *Options) {
		o.Skip = i
	}
}

func WithKeyVals(kv ...interface{}) Option {
	return func(o *Options) {
		o.Keyvals = append(o.Keyvals, kv...)
	}
}

type Logger interface {
	Log(level Level, keyvals ...interface{}) error
	LogWithOptions(opts ...Option) error
}

type logger struct {
	logger Logger
}

func New(l Logger) Logger {
	return &logger{
		logger: l,
	}
}

func (l *logger) Log(level Level, keyvals ...interface{}) error {
	return l.logger.Log(level, keyvals...)
}

func (l *logger) LogWithOptions(opts ...Option) error {
	return l.logger.LogWithOptions(opts...)
}
