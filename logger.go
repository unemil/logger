package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	logger *slog.Logger

	ctxFieldKeys = make(map[any]struct{}, 0)

	logLevels = map[slog.Leveler]string{
		levelTrace: logLevelTrace,
		levelFatal: logLevelFatal,
		levelPanic: logLevelPanic,
	}
)

const (
	logLevel  = "LOG_LEVEL"
	logFields = "LOG_FIELDS"

	logLevelTrace = "TRACE"
	logLevelFatal = "FATAL"
	logLevelPanic = "PANIC"

	levelTrace = slog.Level(slog.LevelDebug - (4 << 0))
	levelFatal = slog.Level(slog.LevelError + (4 << 0))
	levelPanic = slog.Level(slog.LevelError + (4 << 1))
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: func() slog.Level {
			switch os.Getenv(logLevel) {
			case logLevelTrace:
				return levelTrace
			case slog.LevelDebug.String():
				return slog.LevelDebug
			default:
				return slog.LevelInfo
			}
		}(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}

			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)

				pcs := make([]uintptr, 10)
				runtime.Callers(0, pcs)

				frames := runtime.CallersFrames(pcs)
				for {
					frame, more := frames.Next()
					if !more {
						break
					}

					source.File = frame.File
					source.Line = frame.Line
				}

				a.Value = slog.StringValue(fmt.Sprintf("%s:%d",
					source.File[strings.LastIndexByte(source.File[:strings.LastIndexByte(source.File, '/')], '/')+1:],
					source.Line,
				))
			}

			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				value, exists := logLevels[level]
				if !exists {
					value = level.String()
				}

				a.Value = slog.StringValue(value)
			}

			return a
		},
	}))
}

func Info(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg, getAttrs(ctx)...)
}

func Infof(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg, getAttrs(ctx, fields...)...)
}

func Trace(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelTrace, msg, getAttrs(ctx)...)
}

func Tracef(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, levelTrace, msg, getAttrs(ctx, fields...)...)
}

func Debug(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelDebug, msg, getAttrs(ctx)...)
}

func Debugf(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelDebug, msg, getAttrs(ctx, fields...)...)
}

func Warn(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelWarn, msg, getAttrs(ctx)...)
}

func Warnf(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelWarn, msg, getAttrs(ctx, fields...)...)
}

func Error(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, slog.LevelError, msg, getAttrs(ctx, getErrorFields(err)...)...)
}

func Errorf(ctx context.Context, msg string, err error, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelError, msg, getAttrs(ctx, getErrorFields(err, fields...)...)...)
}

func Fatal(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelFatal, msg, getAttrs(ctx, getErrorFields(err)...)...)
	os.Exit(1)
}

func Fatalf(ctx context.Context, msg string, err error, fields ...any) {
	logger.LogAttrs(ctx, levelFatal, msg, getAttrs(ctx, getErrorFields(err, fields...)...)...)
	os.Exit(1)
}

func Panic(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelPanic, msg, getAttrs(ctx, getErrorFields(err)...)...)
	panic(err)
}

func Panicf(ctx context.Context, msg string, err error, fields ...any) {
	logger.LogAttrs(ctx, levelPanic, msg, getAttrs(ctx, getErrorFields(err, fields...)...)...)
	panic(err)
}

func WithContext(ctx context.Context, key, value any) context.Context {
	if _, ok := ctxFieldKeys[key]; !ok {
		ctxFieldKeys[key] = struct{}{}
	}

	return context.WithValue(ctx, key, value)
}

func getErrorFields(err error, fields ...any) []any {
	if err == nil {
		err = errors.New("")
	}

	fields = append(fields, "error", err.Error())
	fields = append(fields[len(fields)-2:], fields[:len(fields)-2]...)

	return fields
}

func getAttrs(ctx context.Context, fields ...any) []slog.Attr {
	ctxFields := make([]any, 0)
	if value := ctx.Value(logFields); value != nil {
		ctxFields = value.([]any)
	}

	if len(fields)%2 != 0 {
		fields = append(fields[:len(fields)-1], fields[len(fields):]...)
	}

	if len(ctxFields)%2 != 0 {
		ctxFields = append(ctxFields[:len(ctxFields)-1], ctxFields[len(ctxFields):]...)
	}

	attrs := make([]slog.Attr, 0, len(fields)+len(ctxFields)+len(ctxFieldKeys))

	for i := 0; i < len(fields); i += 2 {
		attrs = appendAttr(attrs, fields[i], fields[i+1])
	}

	for i := 0; i < len(ctxFields); i += 2 {
		attrs = appendAttr(attrs, ctxFields[i], ctxFields[i+1])
	}

	for key := range ctxFieldKeys {
		if value := ctx.Value(key); value != nil {
			attrs = appendAttr(attrs, key, value)
		}
	}

	return attrs
}

func appendAttr(attrs []slog.Attr, fieldKey, fieldValue any) []slog.Attr {
	var (
		key   = slog.AnyValue(fieldKey).String()
		value any
	)

	v := slog.AnyValue(fieldValue)
	switch v.Kind() {
	case slog.KindBool:
		value = v.Bool()
	case slog.KindDuration:
		value = v.Duration().String()
	case slog.KindFloat64:
		value = v.Float64()
	case slog.KindInt64:
		value = v.Int64()
	case slog.KindString:
		value = v.String()
	case slog.KindTime:
		value = v.Time().Format(time.RFC3339)
	case slog.KindUint64:
		value = v.Uint64()
	default:
		value = v.Any()
	}

	attrs = append(attrs, slog.Any(key, value))

	return attrs
}
