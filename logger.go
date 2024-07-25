package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/unemil/logger/field"
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
func Tracef(ctx context.Context, msg string, fields ...field.Field) {
	logger.LogAttrs(ctx, levelTrace, msg, convertFields(fields...)...)
}

// Debug logs a message at the DEBUG level.
func Debug(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelDebug, msg)
}

// Debugf logs a message at the DEBUG level with additional fields.
func Debugf(ctx context.Context, msg string, fields ...field.Field) {
	logger.LogAttrs(ctx, levelDebug, msg, convertFields(fields...)...)
}

// Info logs a message at the INFO level.
func Info(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelInfo, msg)
}

// Infof logs a message at the INFO level with additional fields.
func Infof(ctx context.Context, msg string, fields ...field.Field) {
	logger.LogAttrs(ctx, levelInfo, msg, convertFields(fields...)...)
}

// Warn logs a message at the WARN level.
func Warn(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelWarn, msg)
}

// Warnf logs a message at the WARN level with additional fields.
func Warnf(ctx context.Context, msg string, fields ...field.Field) {
	logger.LogAttrs(ctx, levelWarn, msg, convertFields(fields...)...)
}

// Error logs a message at the ERROR level with an associated error.
func Error(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelError, msg, convertFields(errorField(err))...)
}

// Errorf logs a message at the ERROR level with an associated error and additional fields.
func Errorf(ctx context.Context, msg string, err error, fields ...field.Field) {
	logger.LogAttrs(ctx, levelError, msg, convertFields(append(fields, errorField(err))...)...)
}

// Fatal logs a message at the FATAL level with an associated error and exits the program.
func Fatal(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelFatal, msg, convertFields(errorField(err))...)
	os.Exit(1)
}

// Fatalf logs a message at the FATAL level with an associated error, additional fields and exits the program.
func Fatalf(ctx context.Context, msg string, err error, fields ...field.Field) {
	logger.LogAttrs(ctx, levelFatal, msg, convertFields(append(fields, errorField(err))...)...)
	os.Exit(1)
}

// Panic logs a message at the PANIC level with an associated error and panics.
func Panic(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelPanic, msg, convertFields(errorField(err))...)
	panic(err)
}

// Panicf logs a message at the PANIC level with an associated error, additional fields and panics.
func Panicf(ctx context.Context, msg string, err error, fields ...field.Field) {
	logger.LogAttrs(ctx, levelPanic, msg, convertFields(append(fields, errorField(err))...)...)
	panic(err)
}

// Field returns a logging field with a specified key-value pair.
func Field(key field.Key, value field.Value) field.Field {
	return field.Field{Key: key, Value: value}
}

// Context returns a context with a specified key-value pair.
func Context(ctx context.Context, key field.Key, value field.Value) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if ctxFields, ok := ctx.Value(fieldsKey).(field.Fields); ok {
		ctxFields[key] = value
		return context.WithValue(ctx, fieldsKey, ctxFields)
	}

	return context.WithValue(ctx, fieldsKey, field.Fields{key: value})
}
