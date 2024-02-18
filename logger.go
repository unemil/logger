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

// Trace logs a message at the TRACE level.
func Trace(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelTrace, msg)
}

// Tracef logs a message at the TRACE level with additional fields.
func Tracef(ctx context.Context, msg string, fields ...Field) {
	logger.LogAttrs(ctx, levelTrace, msg, Fields(fields).toAttrs()...)
}

// Debug logs a message at the DEBUG level.
func Debug(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelDebug, msg)
}

// Debugf logs a message at the DEBUG level with additional fields.
func Debugf(ctx context.Context, msg string, fields ...Field) {
	logger.LogAttrs(ctx, levelDebug, msg, Fields(fields).toAttrs()...)
}

// Info logs a message at the INFO level.
func Info(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelInfo, msg)
}

// Infof logs a message at the INFO level with additional fields.
func Infof(ctx context.Context, msg string, fields ...Field) {
	logger.LogAttrs(ctx, levelInfo, msg, Fields(fields).toAttrs()...)
}

// Warn logs a message at the WARN level.
func Warn(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelWarn, msg)
}

// Warnf logs a message at the WARN level with additional fields.
func Warnf(ctx context.Context, msg string, fields ...Field) {
	logger.LogAttrs(ctx, levelWarn, msg, Fields(fields).toAttrs()...)
}

// Error logs a message at the ERROR level with an associated error.
func Error(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelError, msg, errField(err).toAttr())
}

// Errorf logs a message at the ERROR level with an associated error and additional fields.
func Errorf(ctx context.Context, msg string, err error, fields ...Field) {
	logger.LogAttrs(ctx, levelError, msg, append(Fields{errField(err)}, fields...).toAttrs()...)
}

// Fatal logs a message at the FATAL level with an associated error and exits the program.
func Fatal(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelFatal, msg, errField(err).toAttr())
	os.Exit(1)
}

// Fatalf logs a message at the FATAL level with an associated error, additional fields and exits the program.
func Fatalf(ctx context.Context, msg string, err error, fields ...Field) {
	logger.LogAttrs(ctx, levelFatal, msg, append(Fields{errField(err)}, fields...).toAttrs()...)
	os.Exit(1)
}

// Panic logs a message at the PANIC level with an associated error and panics.
func Panic(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelPanic, msg, errField(err).toAttr())
	panic(err)
}

// Panicf logs a message at the PANIC level with an associated error, additional fields and panics.
func Panicf(ctx context.Context, msg string, err error, fields ...Field) {
	logger.LogAttrs(ctx, levelPanic, msg, append(Fields{errField(err)}, fields...).toAttrs()...)
	panic(err)
}

// With returns a context with a logging key and non-nil value pair.
func With(ctx context.Context, key FieldKey, value FieldValue) context.Context {
	ctxFieldKeys = append(ctxFieldKeys, key)

	return context.WithValue(ctx, key, value)
}
