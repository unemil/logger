package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/unemil/logger/field"
)

type contextHandler struct {
	slog.Handler
}

const (
	levelTrace slog.Level = 4 * (iota - 2)
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
	levelPanic
)

var (
	levels = map[slog.Level]string{
		levelTrace: "TRACE",
		levelDebug: "DEBUG",
		levelInfo:  "INFO",
		levelWarn:  "WARN",
		levelError: "ERROR",
		levelFatal: "FATAL",
		levelPanic: "PANIC",
	}

	ctxFieldKeys = make([]field.Key, 0)

	handlerOptions = slog.HandlerOptions{
		AddSource: true,
		Level: func() slog.Level {
			switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
			case "TRACE":
				return levelTrace
			case "DEBUG":
				return levelDebug
			}

			return levelInfo
		}(),
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
	return &contextHandler{Handler: slog.NewJSONHandler(os.Stdout, &handlerOptions)}
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	var (
		uniqueAttrs = make(map[string]slog.Value, 0)
		attrKeys    = make([]string, 0, len(ctxFieldKeys)+r.NumAttrs())

		removeKey = func(keys []string, key string) []string {
			for i := range keys {
				if keys[i] == key {
					return append(keys[:i], keys[i+1:]...)
				}
			}

			return keys
		}
		setAttr = func(a slog.Attr) {
			if _, ok := uniqueAttrs[a.Key]; ok {
				attrKeys = removeKey(attrKeys, a.Key)
			}
			attrKeys = append(attrKeys, a.Key)
			uniqueAttrs[a.Key] = a.Value
		}
	)

	for _, key := range ctxFieldKeys {
		if value := ctx.Value(key); value != nil {
			setAttr(convertField(field.Field{Key: key, Value: value}))
		}
	}
	r.Attrs(func(a slog.Attr) bool { setAttr(a); return true })

	attrs := make([]slog.Attr, 0, len(attrKeys))
	for _, key := range attrKeys {
		attrs = append(attrs, slog.Attr{Key: key, Value: uniqueAttrs[key]})
	}

	r = slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}

func errorField(err field.Value) field.Field {
	return field.Field{Key: "error", Value: err}
}

func appendFields(f field.Field, fs field.Fields) field.Fields {
	return append(field.Fields{f}, fs...)
}

func convertField(f field.Field) slog.Attr {
	switch v := slog.AnyValue(f.Value); v.Kind() {
	case slog.KindBool:
		f.Value = v.Bool()
	case slog.KindDuration:
		f.Value = v.Duration().String()
	case slog.KindFloat64:
		f.Value = v.Float64()
	case slog.KindInt64:
		f.Value = v.Int64()
	case slog.KindString:
		f.Value = v.String()
	case slog.KindTime:
		f.Value = v.Time().Format(time.RFC3339)
	case slog.KindUint64:
		f.Value = v.Uint64()
	case slog.KindGroup:
		f.Value = v.Group()
	case slog.KindLogValuer:
		f.Value = v.LogValuer()
	default:
		f.Value = v.Any()
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
