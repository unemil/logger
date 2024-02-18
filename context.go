package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

type contextHandler struct {
	slog.Handler
}

var (
	ctxFieldKeys = make([]FieldKey, 0)

	handlerOptions = slog.HandlerOptions{
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
			}

			if a.Key == slog.LevelKey {
				a.Value = slog.AnyValue(levels[a.Value.Any().(slog.Level)])
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
	}
)

func newContextHandler() *contextHandler {
	handlerOptions.Level = logLevel()

	return &contextHandler{Handler: slog.NewJSONHandler(os.Stdout, &handlerOptions)}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		uniqueFields = make(map[FieldKey]FieldValue, 0)
		fieldKeys    = make([]FieldKey, 0, len(ctxFieldKeys)+r.NumAttrs())

		setField = func(key FieldKey, value FieldValue) {
			if _, ok := uniqueFields[key]; !ok {
				fieldKeys = append(fieldKeys, key)
			}
			uniqueFields[key] = value
		}
	)

	for _, key := range ctxFieldKeys {
		if value := ctx.Value(key); value != nil {
			setField(key, value)
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		setField(FieldKey(a.Key), FieldValue(a.Value))

		return true
	})

	attrs := make([]slog.Attr, 0, len(fieldKeys))
	for _, key := range fieldKeys {
		attrs = append(attrs, Field{Key: key, Value: uniqueFields[key]}.toAttr())
	}

	r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
