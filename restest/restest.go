package restest

import "github.com/omcrgnt/res"

// wire registers arbitrary v as a NewResourceer for test setup.
type Wire struct{ V any }

func (w Wire) NewResource() (any, error) { return w.V, nil }


// Registry returns a new empty registry ([res.New]).
func Registry() res.Registry {
	return res.New()
}

// With returns a registry with values registered in order via [res.Registry.Add].
func With(values ...any) res.Registry {
	reg := Registry()
	for _, v := range values {
		Must(reg.Add(Wire{V: v}))
	}
	return reg
}

// WithTagged returns a registry containing one value registered via [res.Registry.AddWithTags].
func WithTagged(v any, tags ...res.Tag) res.Registry {
	reg := Registry()
	Must(reg.AddWithTags(Wire{V: v}, tags...))
	return reg
}

// AddAll registers values in order. Stops and returns the first error.
func AddAll(reg res.Registry, values ...any) error {
	for _, v := range values {
		if err := reg.Add(Wire{V: v}); err != nil {
			return err
		}
	}
	return nil
}

// MustAddAll is [AddAll] that panics on error.
func MustAddAll(reg res.Registry, values ...any) {
	Must(AddAll(reg, values...))
}

// ResetGlobal replaces [res.Global] with an empty registry and returns it.
// Use for tests that exercise the shared registry (e.g. after */use init).
func ResetGlobal() res.Registry {
	return res.ResetGlobalForRestest()
}

// Must panics if err is non-nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
