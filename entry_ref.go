package res

import (
	"fmt"
	"reflect"
)

type entryRef struct {
	reg      *registry
	identity any
}

func (e *entryRef) Type() reflect.Type {
	return reflect.TypeOf(e.identity)
}

func (e *entryRef) Value() any {
	return e.identity
}

func (e *entryRef) item() (item, bool) {
	e.reg.mu.RLock()
	defer e.reg.mu.RUnlock()

	for _, it := range e.reg.items {
		if it.value == e.identity {
			return it, true
		}
	}
	return item{}, false
}

func (e *entryRef) Has(tag Tag) bool {
	it, ok := e.item()
	if !ok {
		return false
	}
	_, ok = it.tags[tag]
	return ok
}

func (e *entryRef) Regular() bool {
	return e.Has(TagRegular)
}

func (e *entryRef) Replaceable() bool {
	return e.Has(TagReplaceable)
}

func (e *entryRef) Fixed() bool {
	return e.Has(TagFixed)
}

func (e *entryRef) Tags() []Tag {
	it, ok := e.item()
	if !ok || len(it.tags) == 0 {
		return nil
	}
	out := make([]Tag, 0, len(it.tags))
	for tag := range it.tags {
		out = append(out, tag)
	}
	return out
}

func (e *entryRef) GetCustomTag(key string) (any, bool) {
	it, ok := e.item()
	if !ok || len(it.customTags) == 0 {
		return nil, false
	}
	val, ok := it.customTags[key]
	return val, ok
}

func (e *entryRef) ChangeValue(new any) error {
	if new == nil {
		return fmt.Errorf("cannot change to nil resource")
	}

	e.reg.mu.Lock()
	defer e.reg.mu.Unlock()

	for i, it := range e.reg.items {
		if it.value == e.identity {
			e.reg.items[i].value = new
			e.identity = new
			return nil
		}
	}
	return fmt.Errorf("resource not found")
}

func newEntryRef(reg *registry, identity any) Entry {
	return &entryRef{reg: reg, identity: identity}
}

// EntryCustomTags returns a copy of custom tags on e, or nil.
func EntryCustomTags(e Entry) map[string]any {
	ref, ok := e.(*entryRef)
	if !ok {
		return nil
	}
	it, ok := ref.item()
	if !ok || len(it.customTags) == 0 {
		return nil
	}
	return cloneCustomTags(it.customTags)
}

func cloneCustomTags(m map[string]any) map[string]any {
	if len(m) == 0 {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
