package field

type (
	// Fields represents a collection of unique logging fields.
	Fields map[Key]Value

	// Field represents a single logging field consisting of a key-value pair.
	Field struct {
		Key   Key
		Value Value
	}

	// FieldKey represents a key used for log fields.
	Key string
	// FieldValue represents a value used for log fields.
	Value any
)
