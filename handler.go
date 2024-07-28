package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/unemil/logger/field"
)

type contextKey string

const (
	levelKey  = "LOG_LEVEL"
	formatKey = "LOG_FORMAT"
	fileKey   = "LOG_FILE"

	fieldsKey contextKey = "LOG_FIELDS"

	stackFramesNumber = 11
)

var handlerOptions = &slog.HandlerOptions{
	AddSource: true,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.TimeKey:
			a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
		case slog.LevelKey:
			a.Value = slog.StringValue(levels[a.Value.Any().(slog.Level)])
		case slog.SourceKey:
			source := a.Value.Any().(*slog.Source)

			pcs := make([]uintptr, stackFramesNumber)
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

		return a
	},
}

type contextHandler struct {
	slog.Handler
}

func newContextHandler() *contextHandler {
	writer := func(file string) io.Writer {
		if len(file) > 0 {
			if f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0600); err == nil {
				return f
			}
		}

		return os.Stdout
	}(os.Getenv(fileKey))

	handlerOptions.Level = func(level string) slog.Level {
		switch level {
		case levelTraceName:
			return levelTrace
		case levelDebugName:
			return levelDebug
		}

		return levelInfo
	}(strings.ToUpper(os.Getenv(levelKey)))

	handler := func(format string) slog.Handler {
		if format == formatJSON {
			return slog.NewJSONHandler(writer, handlerOptions)
		}

		return slog.NewTextHandler(writer, handlerOptions)
	}(strings.ToUpper(os.Getenv(formatKey)))

	return &contextHandler{Handler: handler}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	fields := make(field.Fields)
	if ctxFields, ok := ctx.Value(fieldsKey).(field.Fields); ok {
		for key, value := range ctxFields {
			fields[key] = value
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		fields[field.Key(a.Key)] = field.Value(a.Value)

		return true
	})

	attrs := make([]slog.Attr, 0, len(fields))
	for key, value := range fields {
		switch v := slog.AnyValue(value); v.Kind() {
		case slog.KindTime:
			value = v.Time().Format(time.RFC3339)
		default:
			value = v
		}

		attrs = append(attrs, slog.Any(string(key), value))
	}
	sort.SliceStable(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})

	r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}

func errorField(err field.Value) field.Field {
	return field.Field{Key: "error", Value: err}
}

func convertFields(fields ...field.Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, slog.Any(string(field.Key), field.Value))
	}

	return attrs
}
