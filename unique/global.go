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
