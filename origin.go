package res

import "reflect"

// Origin marks whether a resource was added as a system default or by app wiring.
type Origin int

const (
	System Origin = iota
	User
)

// Entry is a registry item with metadata for SDI cleanup.
type Entry struct {
	Type   reflect.Type
	Value  any
	Origin Origin
}
