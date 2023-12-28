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

var ctxFieldKeys = make(map[FieldKey]struct{}, 0)

type contextHandler struct {
	slog.Handler
}

func newContextHandler() *contextHandler {
	level := levelInfo
	switch os.Getenv(logLevel) {
	case logLevelTrace:
		level = levelTrace
	case logLevelDebug:
		level = levelDebug
	}

	return &contextHandler{slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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

				a.Value = slog.StringValue(fmt.Sprintf("%s:%d",
					source.File[strings.LastIndexByte(source.File[:strings.LastIndexByte(source.File, '/')], '/')+1:],
					source.Line,
				))
			}

			return a
		},
	})}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		ctxFields = make(Fields, 0)
		attrs     = make([]slog.Attr, 0, r.NumAttrs()+len(ctxFieldKeys))
		err       slog.Attr
	)

	for key := range ctxFieldKeys {
		ctxFields[key] = ctx.Value(key)
	}

	r.Attrs(func(a slog.Attr) bool {
		if a.Key == string(errFieldKey) {
			switch r.Level {
			case slog.LevelError, levelFatal, levelPanic:
				err = a
			}
		} else {
			delete(ctxFields, FieldKey(a.Key))
			attrs = append(attrs, a)
		}

		return true
	})

	attrs = append(attrs, ctxFields.toAttrs()...)

	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})

	if !err.Equal(slog.Attr{}) {
		attrs = append([]slog.Attr{err}, attrs...)
	}

	r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
