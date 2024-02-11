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
	ctxFieldKeys = make(map[FieldKey]struct{}, 0)

	ctxHandlerOptions = slog.HandlerOptions{
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

type contextHandler struct {
	slog.Handler
}

func newContextHandler() *contextHandler {
	ctxHandlerOptions.Level = envLogLevel()

	return &contextHandler{
		Handler: slog.NewJSONHandler(os.Stdout, &ctxHandlerOptions),
	}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	ctxFields := make(Fields, 0)
	for key := range ctxFieldKeys {
		if value := ctx.Value(key); value != nil && !key.isError() {
			ctxFields[key] = value
		}
	}

	var (
		err   slog.Attr
		attrs = make([]slog.Attr, 0, r.NumAttrs()+len(ctxFields))
	)
	r.Attrs(func(a slog.Attr) bool {
		switch key := FieldKey(a.Key); key.isError() {
		case true:
			err = a
		default:
			attrs = append(attrs, a)
			delete(ctxFields, key)
		}

		return true
	})

	attrs = append(attrs, ctxFields.toAttrs()...)
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})

	if r.Level >= levelError {
		attrs = append([]slog.Attr{err}, attrs...)
	}

	r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
