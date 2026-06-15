package res

import "reflect"

// Entry wraps one resource in a [Registry].
type Entry struct {
	// Type is the concrete type of Value.
	Type reflect.Type
	// Value is the stored resource.
	Value any
	tags  tagSet
}
