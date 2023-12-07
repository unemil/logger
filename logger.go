package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"
)

var (
	logger = new(slog.Logger)

	ctxFieldKeys   = make(map[FieldKey]struct{}, 0)
	fieldKeysCache = make(map[FieldKey]struct{}, 0)

	logLevels = map[slog.Leveler]string{
		levelTrace: logLevelTrace,
		levelFatal: logLevelFatal,
		levelPanic: logLevelPanic,
	}
)

const (
	logLevel = "LOG_LEVEL"

	logLevelTrace = "TRACE"
	logLevelFatal = "FATAL"
	logLevelPanic = "PANIC"

	levelTrace = slog.Level(slog.LevelDebug - (4 << 0))
	levelFatal = slog.Level(slog.LevelError + (4 << 0))
	levelPanic = slog.Level(slog.LevelError + (4 << 1))

	errFieldKey FieldKey = "error"
)

type (
	FieldKey string

	contextHandler struct{ slog.Handler }
)

func newContextHandler(h slog.Handler) *contextHandler { return &contextHandler{h} }

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0, len(ctxFieldKeys))
	for key := range ctxFieldKeys {
		if value := ctx.Value(key); value != nil {
			attrs = appendAttr(attrs, key, value)
		}
	}

	r.AddAttrs(attrs...)

	var err slog.Attr
	attrs = make([]slog.Attr, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case string(errFieldKey):
			switch r.Level {
			case slog.LevelError, levelFatal, levelPanic:
				err = a
			}
		default:
			attrs = append(attrs, a)
		}

		return true
	})

	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})

	if !err.Equal(slog.Attr{}) {
		attrs = append(attrs, err)
		attrs = append(attrs[len(attrs)-1:], attrs[:len(attrs)-1]...)
	}

	record := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	record.AddAttrs(attrs...)

	fieldKeysCache = make(map[FieldKey]struct{}, 0)

	return h.Handler.Handle(ctx, record)
}

func init() {
	logger = slog.New(newContextHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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

				pcs := make([]uintptr, 11)
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
		}})))
}

func Info(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg)
}

func Infof(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg, getAttrs(nil, fields...)...)
}

func Trace(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, levelTrace, msg)
}

func Tracef(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, levelTrace, msg, getAttrs(nil, fields...)...)
}

func Debug(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelDebug, msg)
}

func Debugf(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelDebug, msg, getAttrs(nil, fields...)...)
}

func Warn(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelWarn, msg)
}

func Warnf(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelWarn, msg, getAttrs(nil, fields...)...)
}

func Error(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, slog.LevelError, msg, getAttrs(err)...)
}

func Errorf(ctx context.Context, msg string, err error, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelError, msg, getAttrs(err, fields...)...)
}

func Fatal(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelFatal, msg, getAttrs(err)...)
	os.Exit(1)
}

func Fatalf(ctx context.Context, msg string, err error, fields ...any) {
	logger.LogAttrs(ctx, levelFatal, msg, getAttrs(err, fields...)...)
	os.Exit(1)
}

func Panic(ctx context.Context, msg string, err error) {
	logger.LogAttrs(ctx, levelPanic, msg, getAttrs(err)...)
	panic(err)
}

func Panicf(ctx context.Context, msg string, err error, fields ...any) {
	logger.LogAttrs(ctx, levelPanic, msg, getAttrs(err, fields...)...)
	panic(err)
}

func WithContext(ctx context.Context, key FieldKey, value any) context.Context {
	ctxFieldKeys[key] = struct{}{}

	return context.WithValue(ctx, key, value)
}

func getAttrs(err error, fields ...any) []slog.Attr {
	if len(fields)%2 != 0 {
		fields = append(fields[:len(fields)-1], fields[len(fields):]...)
	}

	unique := make(map[FieldKey]any, len(fields))
	for i := 0; i < len(fields); i += 2 {
		unique[FieldKey(slog.AnyValue(fields[i]).String())] = fields[i+1]
	}

	attrs := make([]slog.Attr, 0, len(unique))
	attrs = appendAttr(attrs, errFieldKey, err)
	for key, value := range unique {
		attrs = appendAttr(attrs, key, value)
	}

	return attrs
}

func appendAttr(attrs []slog.Attr, fieldKey FieldKey, fieldValue any) []slog.Attr {
	if _, ok := fieldKeysCache[fieldKey]; ok {
		return attrs
	}

	fieldKeysCache[fieldKey] = struct{}{}

	var (
		key   = string(fieldKey)
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
