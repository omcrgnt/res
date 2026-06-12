package res

import (
	"fmt"
	"reflect"
	"sync"
)

var globalRegistry = newRegistry()

// Default — глобальный реестр: builder.Add, sdi.Pool (Walk), Get/Find.
var Default globalAccessor

// Instance — алиас Default.
var Instance = Default

type globalAccessor struct{}

func (globalAccessor) Add(v any) error {
	return addUser(globalRegistry, v)
}

func (globalAccessor) AddBuiltin(v any) error {
	return addBuiltin(globalRegistry, v)
}

func (globalAccessor) Walk(fn func(t reflect.Type, res any) bool) {
	globalRegistry.walk(fn)
}

func (globalAccessor) WalkEntries(fn func(Entry) bool) {
	globalRegistry.walkEntries(fn)
}

func keyOf[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

// AddBuiltin registers a system default resource (from package init).
func AddBuiltin[T any](v T) error {
	return addBuiltin(globalRegistry, v)
}

// Add registers a user resource (from builder.Build).
// Replaces an existing system resource of the same concrete type.
func Add[T any](v T) error {
	return addUser(globalRegistry, v)
}

// AddAll adds user resources to the global registry.
func AddAll(resources ...any) error {
	for _, res := range resources {
		if err := addUser(globalRegistry, res); err != nil {
			return err
		}
	}
	return nil
}

// Remove deletes a resource by value identity.
func Remove(v any) error {
	return removeAny(globalRegistry, v)
}

// Walk performs a read-only walk in registration order.
func Walk(fn func(t reflect.Type, res any) bool) {
	globalRegistry.walk(fn)
}

// WalkEntries walks resources including origin metadata.
func WalkEntries(fn func(Entry) bool) {
	globalRegistry.walkEntries(fn)
}

// Get returns a resource by concrete type.
func Get[T any]() (T, bool) {
	return get[T](globalRegistry)
}

func get[T any](r *registry) (T, bool) {
	var zero T
	key := keyOf[T]()

	r.mu.RLock()
	defer r.mu.RUnlock()

	e, ok := r.byType[key]
	if !ok {
		return zero, false
	}

	return e.value.(T), true
}

// Find returns resources matching T (interface or concrete).
func Find[T any]() []T {
	var matches []T
	targetType := reflect.TypeFor[T]()

	globalRegistry.walkEntries(func(e Entry) bool {
		switch targetType.Kind() {
		case reflect.Interface:
			if e.Type.Implements(targetType) {
				if val, ok := e.Value.(T); ok {
					matches = append(matches, val)
				}
			}
		default:
			if e.Type == targetType {
				if val, ok := e.Value.(T); ok {
					matches = append(matches, val)
				}
			}
		}
		return true
	})

	return matches
}

type entry struct {
	value  any
	origin Origin
}

type registry struct {
	mu        sync.RWMutex
	resources []entry
	byType    map[reflect.Type]entry
}

func newRegistry() *registry {
	return &registry{
		byType: make(map[reflect.Type]entry),
	}
}

func addBuiltin(r *registry, res any) error {
	if res == nil {
		return fmt.Errorf("cannot add nil resource")
	}

	key := reflect.TypeOf(res)

	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.byType[key]; ok {
		if existing.origin == System {
			return fmt.Errorf("resource of type %v is already registered", key)
		}
		return fmt.Errorf("resource of type %v is already registered by user", key)
	}

	e := entry{value: res, origin: System}
	r.byType[key] = e
	r.resources = append(r.resources, e)
	return nil
}

func addUser(r *registry, res any) error {
	if res == nil {
		return fmt.Errorf("cannot add nil resource")
	}

	key := reflect.TypeOf(res)

	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.byType[key]
	if !ok {
		e := entry{value: res, origin: User}
		r.byType[key] = e
		r.resources = append(r.resources, e)
		return nil
	}

	if existing.origin == User {
		return fmt.Errorf("resource of type %v is already registered", key)
	}

	e := entry{value: res, origin: User}
	r.byType[key] = e
	for i := range r.resources {
		if r.resources[i].value == existing.value {
			r.resources[i] = e
			return nil
		}
	}
	return fmt.Errorf("resource of type %v: internal inconsistency", key)
}

func removeAny(r *registry, v any) error {
	if v == nil {
		return fmt.Errorf("cannot remove nil resource")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	idx := -1
	var key reflect.Type
	for i, e := range r.resources {
		if e.value == v {
			idx = i
			key = reflect.TypeOf(v)
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("resource not found")
	}

	delete(r.byType, key)
	r.resources = append(r.resources[:idx], r.resources[idx+1:]...)
	return nil
}

func (r *registry) walk(fn func(t reflect.Type, res any) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.resources {
		if !fn(reflect.TypeOf(e.value), e.value) {
			break
		}
	}
}

func (r *registry) walkEntries(fn func(Entry) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, e := range r.resources {
		if !fn(Entry{
			Type:   reflect.TypeOf(e.value),
			Value:  e.value,
			Origin: e.origin,
		}) {
			break
		}
	}
}

func resetGlobalRegistry() {
	globalRegistry = newRegistry()
}
