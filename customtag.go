package res

import (
	"fmt"
	"reflect"
)

// AddWithTagsAndCustomTags registers v with system tags and custom tags.
// Intended for github.com/omcrgnt/res/unique.
func AddWithTagsAndCustomTags(r Registry, v any, customTags map[string]any, tags ...Tag) error {
	if v == nil {
		return fmt.Errorf("cannot add nil resource")
	}
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag required")
	}

	reg, ok := r.(*registry)
	if !ok {
		return fmt.Errorf("res: AddWithTagsAndCustomTags: unsupported registry implementation")
	}

	reg.mu.Lock()
	defer reg.mu.Unlock()

	reg.items = append(reg.items, item{
		value:      v,
		tags:       newTagSet(tags...),
		customTags: cloneCustomTags(customTags),
	})
	return nil
}

// ReplaceAtTypeWithCustomTags replaces the first entry of v's type, preserving custom tags
// from the replaced entry when the new entry does not supply custom tags.
func ReplaceAtTypeWithCustomTags(r Registry, v any, customTags map[string]any, tags ...Tag) error {
	if v == nil {
		return fmt.Errorf("cannot replace nil resource")
	}
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag required")
	}

	reg, ok := r.(*registry)
	if !ok {
		return fmt.Errorf("res: ReplaceAtTypeWithCustomTags: unsupported registry implementation")
	}

	typ := reflect.TypeOf(v)
	ts := newTagSet(tags...)

	reg.mu.Lock()
	defer reg.mu.Unlock()

	for i, it := range reg.items {
		if reflect.TypeOf(it.value) == typ {
			ct := cloneCustomTags(customTags)
			if ct == nil {
				ct = cloneCustomTags(it.customTags)
			}
			reg.items[i] = item{value: v, tags: ts, customTags: ct}
			return nil
		}
	}

	reg.items = append(reg.items, item{
		value:      v,
		tags:       ts,
		customTags: cloneCustomTags(customTags),
	})
	return nil
}
