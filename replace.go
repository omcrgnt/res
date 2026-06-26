package res

import (
	"fmt"
	"reflect"
)

// ReplaceAtType replaces the first entry whose type equals reflect.TypeOf(v),
// preserving registration order. If none exists, appends like Add or AddWithTags.
// Storage-only: does not interpret tags.
func ReplaceAtType(r Registry, v any, tags ...Tag) error {
	if v == nil {
		return fmt.Errorf("cannot replace nil resource")
	}

	reg, ok := r.(*registry)
	if !ok {
		return fmt.Errorf("res: ReplaceAtType: unsupported registry implementation")
	}

	var ts tagSet
	if len(tags) > 0 {
		ts = newTagSet(tags...)
	}

	typ := reflect.TypeOf(v)

	reg.mu.Lock()
	defer reg.mu.Unlock()

	for i, it := range reg.items {
		if reflect.TypeOf(it.value) == typ {
			reg.items[i] = item{value: v, tags: ts}
			return nil
		}
	}

	reg.items = append(reg.items, item{value: v, tags: ts})
	return nil
}
