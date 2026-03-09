package utilities

import (
	"context"
	"os"
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

type Logger struct{ *zap.Logger }

var global atomic.Value

func Global() *Logger {
	if lg := global.Load(); lg != nil {
		return lg.(*Logger)
	}
	l, err := NewProduction()
	if err != nil {
		l = NewConsole(InfoLevel, true)
	}
	global.Store(l)
	return l
}

func SetGlobal(l *Logger) {
	if l != nil {
		global.Store(l)
	}
}

func buildFrom(cfg zap.Config) (*Logger, error) {
	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return new(Logger{z}), nil
}

func NewDevelopment() (*Logger, error) { return buildFrom(zap.NewDevelopmentConfig()) }
func NewProduction() (*Logger, error)  { return buildFrom(zap.NewProductionConfig()) }

func NewConsole(level zapcore.Level, json bool) *Logger {
	enc := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	if json {
		enc = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zapcore.Level(level)))
	return new(Logger{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))})
}

func (l *Logger) with(fields ...zap.Field) *Logger { return new(Logger{l.Logger.With(fields...)}) }

func fromContext(ctx context.Context) *Logger {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value("logger"); v != nil {
		if l, ok := v.(*Logger); ok {
			return l
		}
	}
	return nil
}

func requestID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	if v := ctx.Value("request_id"); v != nil {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func loggerFromContext(ctx context.Context) *Logger {
	l := fromContext(ctx)
	if l == nil {
		l = Global()
	}
	if id, ok := requestID(ctx); ok && id != "" {
		return l.with(zap.String(string("request_id"), id))
	}
	return l
}

func log(ctx context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	if l := loggerFromContext(ctx); l != nil {
		switch level {
		case DebugLevel:
			l.Debug(msg, fields...)
		case InfoLevel:
			l.Info(msg, fields...)
		case WarnLevel:
			l.Warn(msg, fields...)
		case ErrorLevel:
			l.Error(msg, fields...)
		case FatalLevel:
			l.Fatal(msg, fields...)
		default:
			l.Info(msg, fields...)
		}
	}
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	log(ctx, DebugLevel, msg, fields...)
}
func Info(ctx context.Context, msg string, fields ...zap.Field) { log(ctx, InfoLevel, msg, fields...) }
func Warn(ctx context.Context, msg string, fields ...zap.Field) { log(ctx, WarnLevel, msg, fields...) }
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	log(ctx, ErrorLevel, msg, fields...)
}
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	log(ctx, FatalLevel, msg, fields...)
}
