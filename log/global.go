package log

import (
	"fmt"
	"os"
	"sync"
)

var DefaultMessageKey = "msg"

var global = &globalLogger{}

func init() {
	global.SetLogger(DefaultLogger)
}

type globalLogger struct {
	lock sync.Mutex
	Logger
}

func (g *globalLogger) SetLogger(logger Logger) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.Logger = logger
}

func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

func LogWithOptions(opts ...Option) {
	_ = global.LogWithOptions(opts...)
}

func Log(level Level, keyvals ...interface{}) {
	_ = global.Log(level, keyvals...)
}

func Debug(a ...interface{}) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprint(a...))
}

func Debugf(format string, a ...interface{}) {
	_ = global.Log(LevelDebug, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Debugw(keyvals ...interface{}) {
	_ = global.Log(LevelDebug, keyvals...)
}

func Info(a ...interface{}) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprint(a...))
}

func Infof(format string, a ...interface{}) {
	_ = global.Log(LevelInfo, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Infow(keyvals ...interface{}) {
	_ = global.Log(LevelInfo, keyvals...)
}

func Warn(a ...interface{}) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprint(a...))
}

func Warnf(format string, a ...interface{}) {
	_ = global.Log(LevelWarn, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Warnw(keyvals ...interface{}) {
	_ = global.Log(LevelWarn, keyvals...)
}

func Error(a ...interface{}) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprint(a...))
}

func Errorf(format string, a ...interface{}) {
	_ = global.Log(LevelError, DefaultMessageKey, fmt.Sprintf(format, a...))
}

func Errorw(keyvals ...interface{}) {
	_ = global.Log(LevelError, keyvals...)
}

func Fatal(a ...interface{}) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprint(a...))
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	_ = global.Log(LevelFatal, DefaultMessageKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func Fatalw(keyvals ...interface{}) {
	_ = global.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
