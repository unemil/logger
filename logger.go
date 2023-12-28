package logger

import (
	"context"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(newContextHandler())
}

func Trace(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelTrace, msg)
}

func Tracef(ctx context.Context, msg string, fields Fields) {
	logger.LogAttrs(ctx, levelTrace, msg, fields.toAttrs()...)
}

func Debug(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelDebug, msg)
}

func Debugf(ctx context.Context, msg string, fields Fields) {
	logger.LogAttrs(ctx, slog.LevelDebug, msg, fields.toAttrs()...)
}

func Info(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg)
}

func Infof(ctx context.Context, msg string, fields Fields) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg, fields.toAttrs()...)
}

func Warn(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelWarn, msg)
}

func Warnf(ctx context.Context, msg string, fields Fields) {
	logger.LogAttrs(ctx, slog.LevelWarn, msg, fields.toAttrs()...)
}

func Error(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, slog.LevelError, msg, Fields{errFieldKey: err}.toAttrs()...)
}

func Errorf(ctx context.Context, msg string, err error, fields Fields) {
	logger.LogAttrs(ctx, slog.LevelError, msg, append(Fields{errFieldKey: err}.toAttrs(), fields.toAttrs()...)...)
}

func Fatal(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelFatal, msg, Fields{errFieldKey: err}.toAttrs()...)
	os.Exit(1)
}

func Fatalf(ctx context.Context, msg string, err error, fields Fields) {
	logger.LogAttrs(ctx, levelFatal, msg, append(Fields{errFieldKey: err}.toAttrs(), fields.toAttrs()...)...)
	os.Exit(1)
}

func Panic(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelPanic, msg, Fields{errFieldKey: err}.toAttrs()...)
	panic(err)
}

func Panicf(ctx context.Context, msg string, err error, fields Fields) {
	logger.LogAttrs(ctx, levelPanic, msg, append(Fields{errFieldKey: err}.toAttrs(), fields.toAttrs()...)...)
	panic(err)
}

func With(ctx context.Context, key FieldKey, value FieldValue) context.Context {
	ctxFieldKeys[key] = struct{}{}

	return context.WithValue(ctx, key, value)
}
