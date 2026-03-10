package logger

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

var (
	global atomic.Value
	ctxKey struct{}
)

func Global() *zap.Logger {
	if lg := global.Load(); lg != nil {
		return lg.(*zap.Logger)
	}
	l, err := zap.NewProduction()
	if err != nil {
		l = NewConsole(InfoLevel, true)
	}
	global.Store(l)
	return l
}

func SetGlobal(l *zap.Logger) {
	if l != nil {
		global.Store(l)
	}
}

func NewConsole(level zapcore.Level, json bool) *zap.Logger {
	encCfg := zap.NewProductionEncoderConfig()
	var enc zapcore.Encoder
	if json {
		enc = zapcore.NewJSONEncoder(encCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encCfg)
	}
	return zap.New(zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(level)), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func NewProduction() (*zap.Logger, error) {
	return zap.NewProduction()
}

func NewDevelopment() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func Log(ctx context.Context, level zapcore.Level, msg string, fields ...zap.Field) {
	if ctx != nil {
		if v := ctx.Value(ctxKey); v != nil {
			if f, ok := v.([]zap.Field); ok {
				fields = append(f, fields...)
			}
		}
	}
	Global().Log(level, msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Log(ctx, DebugLevel, msg, fields...)
}
func Info(ctx context.Context, msg string, fields ...zap.Field) { Log(ctx, InfoLevel, msg, fields...) }
func Warn(ctx context.Context, msg string, fields ...zap.Field) { Log(ctx, WarnLevel, msg, fields...) }
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Log(ctx, ErrorLevel, msg, fields...)
}
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	Log(ctx, FatalLevel, msg, fields...)
}
