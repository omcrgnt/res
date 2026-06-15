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
	// GetByType returns entries whose [Entry.Type] equals t.
	GetByType(t reflect.Type) []Entry
	// GetByInterface returns entries whose [Entry.Type] implements iface.
	GetByInterface(iface reflect.Type) []Entry
	// GetOneByType returns the first [Entry.Value] for t in registration order.
	// See also [Registry.GetByType].
	GetOneByType(t reflect.Type) (any, error)
	// GetOneByInterface returns the first matching [Entry.Value] in registration order.
	// See also [Registry.GetByInterface].
	GetOneByInterface(iface reflect.Type) (any, error)

	// Transform applies [TransformFunc] to every stored resource in place.
	Transform(...TransformFunc) error

	// Remove unregisters a resource by Value identity (==).
	Remove(v any) error
}

// Default is the global [Registry].
var Default Registry = New()

// Instance is an alias for [Default].
var Instance = Default

// New returns an empty [Registry].
func New() Registry {
	return &registry{}
}

// Add appends to [Default] via [Registry.Add].
func Add[T any](v T) error {
	return Default.Add(v)
}

// AddWithTags appends to [Default] via [Registry.AddWithTags].
func AddWithTags[T any](v T, tags ...Tag) error {
	return Default.AddWithTags(v, tags...)
}

// WalkEntries walks [Default] via [Registry.WalkEntries].
func WalkEntries(fn func(Entry) bool) {
	Default.WalkEntries(fn)
}

// GetByType queries [Default] via [Registry.GetByType].
func GetByType(t reflect.Type) []Entry {
	return Default.GetByType(t)
}

// GetByInterface queries [Default] via [Registry.GetByInterface].
func GetByInterface(iface reflect.Type) []Entry {
	return Default.GetByInterface(iface)
}

// GetOneByType returns the first match from [Default] via [Registry.GetOneByType].
func GetOneByType(t reflect.Type) (any, error) {
	return Default.GetOneByType(t)
}

// GetOneByInterface returns the first match from [Default] via [Registry.GetOneByInterface].
func GetOneByInterface(iface reflect.Type) (any, error) {
	return Default.GetOneByInterface(iface)
}

// Transform transforms all entries in [Default] via [Registry.Transform].
func Transform(fns ...TransformFunc) error {
	return Default.Transform(fns...)
}

// Remove deletes from [Default] via [Registry.Remove].
func Remove(v any) error {
	return Default.Remove(v)
}

type item struct {
	value any
	tags  tagSet
}

type registry struct {
	mu    sync.RWMutex
	items []item
}

var _ Registry = (*registry)(nil)

func (r *registry) Add(v any) error {
	return r.add(v, nil)
}

func (r *registry) AddWithTags(v any, tags ...Tag) error {
	if len(tags) == 0 {
		return fmt.Errorf("at least one tag required")
	}
	return r.add(v, newTagSet(tags...))
}

func (r *registry) add(v any, tags tagSet) error {
	if v == nil {
		return fmt.Errorf("cannot add nil resource")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = append(r.items, item{value: v, tags: tags})
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
	defer r.mu.RUnlock()

	for _, it := range r.items {
		if !fn(entryFromItem(it)) {
			break
		}
	}
}

func (r *registry) GetByType(t reflect.Type) []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []Entry
	for _, it := range r.items {
		e := entryFromItem(it)
		if e.Type == t {
			out = append(out, e)
		}
	}
	return out
}

func (r *registry) GetByInterface(iface reflect.Type) []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []Entry
	for _, it := range r.items {
		e := entryFromItem(it)
		if e.Type.Implements(iface) {
			out = append(out, e)
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

func entryFromItem(it item) Entry {
	return Entry{
		Type:  reflect.TypeOf(it.value),
		Value: it.value,
		tags:  it.tags,
	}
}

func resetGlobalRegistry() {
	Default = New()
	Instance = Default
}
