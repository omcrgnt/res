package res

import "reflect"

// Entry is a handle to one resource slot in a [Registry].
type Entry interface {
	Type() reflect.Type
	Value() any

	Has(tag Tag) bool
	Regular() bool
	Replaceable() bool
	Fixed() bool
	Tags() []Tag

	GetCustomTag(key string) (any, bool)
	ChangeValue(new any) error
}
