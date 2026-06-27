package res

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrNotFound is returned by [Registry.GetOneByType] and [Registry.GetOneByInterface]
// when no matching entry exists.
var ErrNotFound = errors.New("resource not found")

// Registry stores application resources as [Entry] values.
type Registry interface {
	// Add registers a resource without tags.
	Add(v any) error
	// AddWithTags registers a resource with the given tags (duplicates ignored).
	AddWithTags(v any, tags ...Tag) error

	// WalkEntries visits entries in registration order.
	WalkEntries(fn func(Entry) bool)
	// GetByType returns entries whose concrete type equals t.
	GetByType(t reflect.Type) []Entry
	// GetByInterface returns entries whose concrete type implements iface.
	GetByInterface(iface reflect.Type) []Entry
	// GetOneByType returns the first [Entry.Value] for t in registration order.
	GetOneByType(t reflect.Type) (any, error)
	// GetOneByInterface returns the first matching [Entry.Value] in registration order.
	GetOneByInterface(iface reflect.Type) (any, error)

	// Transform applies [TransformFunc] to every stored resource in place.
	Transform(...TransformFunc) error

	// Remove unregisters a resource by Value identity (==).
	Remove(v any) error
}

var global Registry = New()

// Global returns the shared application [Registry] populated by library use init
// (via [AddToGlobalWithTags]) and used as the composition-root registry.
func Global() Registry {
	return global
}

// AddToGlobalWithTags is [Registry.AddWithTags] on [Global].
func AddToGlobalWithTags(v any, tags ...Tag) error {
	return global.AddWithTags(v, tags...)
}

// New returns an empty [Registry].
func New() Registry {
	return &registry{}
}

type item struct {
	value      any
	tags       tagSet
	customTags map[string]any
}

type registry struct {
	mu    sync.RWMutex
	items []item
}

var _ Registry = (*registry)(nil)

func (r *registry) Add(v any) error {
	return r.add(v, nil, nil)
}

func (r *registry) AddWithTags(v any, tags ...Tag) error {
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag required")
	}
	return r.add(v, newTagSet(tags...), nil)
}

func (r *registry) add(v any, tags tagSet, customTags map[string]any) error {
	if v == nil {
		return fmt.Errorf("cannot add nil resource")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = append(r.items, item{
		value:      v,
		tags:       tags,
		customTags: cloneCustomTags(customTags),
	})
	return nil
}

func (r *registry) Remove(v any) error {
	if v == nil {
		return fmt.Errorf("cannot remove nil resource")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, it := range r.items {
		if it.value == v {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("resource not found")
}

func (r *registry) WalkEntries(fn func(Entry) bool) {
	r.mu.RLock()
	refs := make([]Entry, 0, len(r.items))
	for _, it := range r.items {
		refs = append(refs, newEntryRef(r, it.value))
	}
	r.mu.RUnlock()

	for _, ref := range refs {
		if !fn(ref) {
			break
		}
	}
}

func (r *registry) GetByType(t reflect.Type) []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []Entry
	for _, it := range r.items {
		if reflect.TypeOf(it.value) == t {
			out = append(out, newEntryRef(r, it.value))
		}
	}
	return out
}

func (r *registry) GetByInterface(iface reflect.Type) []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []Entry
	for _, it := range r.items {
		if reflect.TypeOf(it.value).Implements(iface) {
			out = append(out, newEntryRef(r, it.value))
		}
	}
	return out
}

func (r *registry) GetOneByType(t reflect.Type) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, it := range r.items {
		if reflect.TypeOf(it.value) == t {
			return it.value, nil
		}
	}
	return nil, ErrNotFound
}

func (r *registry) GetOneByInterface(iface reflect.Type) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, it := range r.items {
		if reflect.TypeOf(it.value).Implements(iface) {
			return it.value, nil
		}
	}
	return nil, ErrNotFound
}

func resetGlobalRegistry() {
	global = New()
}
