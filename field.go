package logger

import (
	"log/slog"
	"time"
)

type (
	// Fields represents a collection of logging fields.
	Fields []Field

	// Field represents a single logging field consisting of a key-value pair.
	Field struct {
		Key   FieldKey
		Value FieldValue
	}

	// FieldKey represents a key used for log fields.
	FieldKey string
	// FieldValue represents a value used for log fields.
	FieldValue any
)

const errFieldKey FieldKey = "error"

func errField(err FieldValue) Field {
	return Field{Key: errFieldKey, Value: err}
}

func (f Field) toAttr() slog.Attr {
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

func (fs Fields) toAttrs() []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fs))
	for _, f := range fs {
		attrs = append(attrs, f.toAttr())
	}

	return attrs
}
