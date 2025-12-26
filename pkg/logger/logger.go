package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type CtxKey string

const LoggerKey CtxKey = "logger"

var log *zap.Logger

func Init() {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewJSONEncoder(config)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	log = zap.New(core, zap.AddCaller())
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return l
	}
	return log
}

func WithContext(ctx context.Context, fields ...zap.Field) context.Context {
	l := log.With(fields...)
	return context.WithValue(ctx, LoggerKey, l)
}
