package unique

import "github.com/omcrgnt/res"

var global = New()

// Global returns the shared type-unique composition-root registry.
func Global() *Registry {
	return global
}

// New returns an empty type-unique registry backed by [res.New].
func New() *Registry {
	return newRegistry(res.New())
}

// MustAddReplaceable registers v on [Global] with [res.TagReplaceable]; panics on error.
func MustAddReplaceable(v any) {
	if err := Global().AddReplaceable(v); err != nil {
		panic(err)
	}
}

// MustAddFixed registers v on [Global] with [res.TagFixed]; panics on error.
func MustAddFixed(v any) {
	if err := Global().AddFixed(v); err != nil {
		panic(err)
	}
}
