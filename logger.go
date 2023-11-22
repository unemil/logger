package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var logger *slog.Logger

const LogLevel = "LOG_LEVEL" // TODO: rename or set custom ???

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level: func() slog.Level {
			logLevelEnv := os.Getenv(LogLevel)
			switch logLevelEnv {
			case slog.LevelDebug.String():
				return slog.LevelDebug
			case slog.LevelWarn.String():
				return slog.LevelWarn
			case slog.LevelError.String():
				return slog.LevelError
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
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", source.File, source.Line)) // TODO: maybe a relative path + line ???
			}

			return a
		},
	}))
}

// TODO: add tests

func Info(ctx context.Context, msg string) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg, getAttrs(ctx)...)
}

func Infof(ctx context.Context, msg string, fields ...any) {
	logger.LogAttrs(ctx, slog.LevelInfo, msg, getAttrs(ctx, fields...)...)
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
	if err == nil {
		err = errors.New("")
	}

	logger.LogAttrs(ctx, slog.LevelError, msg, getAttrs(ctx, "error", err.Error())...)
}

func Errorf(ctx context.Context, msg string, err error, fields ...any) {
	if err == nil {
		err = errors.New("")
	}

	fields = append(fields, "error", err.Error())
	fields = append(fields[len(fields)-2:], fields[:len(fields)-2]...)

	logger.LogAttrs(ctx, slog.LevelError, msg, getAttrs(ctx, fields...)...)
}

func getAttrs(ctx context.Context, fields ...any) []slog.Attr {
	if len(fields)%2 != 0 {
		fields = append(fields[:len(fields)-1], fields[len(fields):]...) // TODO: delete or append ???
	}

	// TODO: get fields from context

	attrs := make([]slog.Attr, 0, len(fields))
	for i := 0; i < len(fields); i += 2 {
		var (
			key   = slog.AnyValue(fields[i]).String()
			value any
		)

		v := slog.AnyValue(fields[i+1])
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
	}

	return attrs
}
