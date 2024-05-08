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

	"github.com/unemil/logger/field"
)

const (
	levelKey  = "LOG_LEVEL"
	fieldsKey = "LOG_FIELDS"
)

type contextHandler struct {
	slog.Handler
}

func newContextHandler() *contextHandler {
	level := levelInfo
	switch strings.ToUpper(os.Getenv(levelKey)) {
	case levelTraceName:
		level = levelTrace
	case levelDebugName:
		level = levelDebug
	}

	return &contextHandler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
				}

				if a.Key == slog.LevelKey {
					a.Value = slog.StringValue(levels[a.Value.Any().(slog.Level)])
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

					a.Value = slog.StringValue(fmt.Sprintf(
						"%s:%d",
						source.File[strings.LastIndexByte(source.File[:strings.LastIndexByte(source.File, '/')], '/')+1:],
						source.Line,
					))
				}

				return a
			},
		}),
	}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		uniqueAttrs = make(map[string]any)
		setAttr     = func(a slog.Attr) {
			uniqueAttrs[a.Key] = a.Value
		}
	)

	if attrs, ok := ctx.Value(fieldsKey).([]slog.Attr); ok {
		for _, attr := range attrs {
			setAttr(attr)
		}
	}
	r.Attrs(func(a slog.Attr) bool { setAttr(a); return true })

	attrs := make([]slog.Attr, 0, len(uniqueAttrs))
	for key, value := range uniqueAttrs {
		attrs = append(attrs, slog.Any(key, value))
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

func convertField(f field.Field) slog.Attr {
	switch v := slog.AnyValue(f.Value); v.Kind() {
	case slog.KindTime:
		f.Value = v.Time().Format(time.RFC3339)
	default:
		f.Value = v
	}

	return slog.Any(string(f.Key), f.Value)
}

func convertFields(fs field.Fields) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fs))
	for _, f := range fs {
		attrs = append(attrs, convertField(f))
	}

	return attrs
}
