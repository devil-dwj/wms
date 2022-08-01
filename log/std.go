package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
)

type stdLogger struct {
	log     *log.Logger
	bufPool *sync.Pool
}

func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		log: log.New(w, "", 0),
		bufPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (l *stdLogger) Log(level Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) & 1) == 1 {
		keyvals = append(keyvals, "KV should be paired")
	}
	buf := l.bufPool.Get().(*bytes.Buffer)
	buf.WriteString(level.String())
	for i := 0; i < len(keyvals); i += 2 {
		_, _ = fmt.Fprintf(buf, " %s=%v", keyvals[i], keyvals[i+1])
	}
	_ = l.log.Output(4, buf.String())
	buf.Reset()
	l.bufPool.Put(buf)
	return nil
}

func (l *stdLogger) LogWithOptions(opts ...Option) error {
	o := Options{
		Level: LevelInfo,
		Skip:  2,
	}
	for _, opt := range opts {
		opt(&o)
	}
	l.Log(o.Level, o.Keyvals...)
	return nil
}
