package res

import "fmt"

// ReplaceEntryForBuilder swaps old for new in reg, preserving entry tags.
// Used by builder.Build to materialize wire entries in place.
//
//go:internal github.com/omcrgnt/builder
func ReplaceEntryForBuilder(reg Registry, old, new any) error {
	r, ok := reg.(*registry)
	if !ok {
		return fmt.Errorf("res: replace: unsupported registry type %T", reg)
	}
	return r.replaceValue(old, new)
}

func (r *registry) replaceValue(old, new any) error {
	if old == nil || new == nil {
		return fmt.Errorf("cannot replace nil resource")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, it := range r.items {
		if it.value == old {
			r.items[i].value = new
			return nil
		}
	}
	return fmt.Errorf("resource not found")
}
