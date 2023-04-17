package zap

import (
	"fmt"
	"os"

	"github.com/devil-dwj/wms/log"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger() log.Logger {
	return &Logger{
		log: MustLog("logs/log.log"),
	}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("KV should be paired: ", keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), fmt.Sprint(keyvals[i+1])))
	}
	switch level {
	case log.LevelDebug:
		l.log.Debug("", data...)
	case log.LevelInfo:
		l.log.Info("", data...)
	case log.LevelWarn:
		l.log.Warn("", data...)
	case log.LevelError:
		l.log.Error("", data...)
	case log.LevelFatal:
		l.log.Fatal("", data...)
	}
	return nil
}

func (l *Logger) LogWithOptions(opts ...log.Option) error {
	o := log.Options{}
	for _, opt := range opts {
		opt(&o)
	}
	if len(o.Keyvals) == 0 || len(o.Keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("KV should be paired: ", o.Keyvals))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(o.Keyvals); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(o.Keyvals[i]), fmt.Sprint(o.Keyvals[i+1])))
	}
	switch o.Level {
	case log.LevelDebug:
		l.log.WithOptions(zap.AddCallerSkip(o.Skip)).Debug("", data...)
	case log.LevelInfo:
		l.log.WithOptions(zap.AddCallerSkip(o.Skip)).Info("", data...)
	case log.LevelWarn:
		l.log.WithOptions(zap.AddCallerSkip(o.Skip)).Warn("", data...)
	case log.LevelError:
		l.log.WithOptions(zap.AddCallerSkip(o.Skip)).Error("", data...)
	case log.LevelFatal:
		l.log.WithOptions(zap.AddCallerSkip(o.Skip)).Fatal("", data...)
	}

	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Close() error {
	return l.Sync()
}

func MustLog(logname string) *zap.Logger {
	w := getWriter(logname)
	muw := zapcore.NewMultiWriteSyncer(w, os.Stdout)
	e := getEncoder()

	c := zapcore.NewCore(e, muw, zapcore.DebugLevel)
	l := zap.New(c, zap.AddCaller(), zap.AddCallerSkip(2))

	zap.ReplaceGlobals(l)

	return l
}

func getWriter(name string) zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   name,
		MaxSize:    100,
		MaxBackups: 500,
		MaxAge:     30,
		LocalTime:  true,
		Compress:   false,
	}
	return zapcore.AddSync(l)
}

func getEncoder() zapcore.Encoder {
	e := zap.NewProductionEncoderConfig()
	e.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	e.TimeKey = "time"
	e.EncodeLevel = zapcore.CapitalLevelEncoder
	e.EncodeDuration = zapcore.MillisDurationEncoder
	e.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewConsoleEncoder(e)
}
