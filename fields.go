package logger

import (
	"log/slog"
	"time"
)

type (
	// Fields represents a collection of fields for logging
	Fields map[FieldKey]FieldValue

	// FieldKey represents a key used for log fields
	FieldKey string
	// FieldValue represents a value used for log fields
	FieldValue any
)

const errFieldKey FieldKey = "error"

func (fs Fields) toAttrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fs))
	for key, value := range fs {
		v := slog.AnyValue(value)
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
		case slog.KindGroup:
			value = v.Group()
		case slog.KindLogValuer:
			value = v.LogValuer()
		default:
			value = v.Any()
		}

		attrs = append(attrs, slog.Any(string(key), value))
	}

	return attrs
}
