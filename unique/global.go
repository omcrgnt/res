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
	Global().MustAddReplaceable(v)
}

// MustAddFixed registers v on [Global] with [res.TagFixed]; panics on error.
func MustAddFixed(v any) {
	Global().MustAddFixed(v)
}

// MustAddReplaceable registers v with [res.TagReplaceable]; panics on error.
func (r *Registry) MustAddReplaceable(v any) {
	if err := r.addReplaceable(v); err != nil {
		panic(err)
	}
}

// MustAddFixed registers v with [res.TagFixed]; panics on error.
func (r *Registry) MustAddFixed(v any) {
	if err := r.addFixed(v); err != nil {
		panic(err)
	}
}
