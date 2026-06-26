package respolicy

import "fmt"

// AddPolicy runs before each Add or AddWithTags on a wrapped [Registry].
type AddPolicy interface {
	PrepareAdd(v any) (any, error)
}

// AcceptAll stores v unchanged.
type AcceptAll struct{}

func (AcceptAll) PrepareAdd(v any) (any, error) {
	return v, nil
}

// RejectInvalid rejects nil values. Further validation is TODO.
type RejectInvalid struct{}

func (RejectInvalid) PrepareAdd(v any) (any, error) {
	if v == nil {
		return nil, fmt.Errorf("respolicy: cannot add nil resource")
	}
	return v, nil
}

func applyPolicies(v any, policies []AddPolicy) (any, error) {
	for _, p := range policies {
		if p == nil {
			continue
		}
		var err error
		v, err = p.PrepareAdd(v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}
