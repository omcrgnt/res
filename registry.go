package res

import (
	"fmt"
	"reflect"
	"sync"
)

var globalRegistry = newRegistry()

func keyOf[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func Add[T any](t T) error {
	return add(globalRegistry, t)
}

func add[T any](r *registry, t T) error {
	key := keyOf[T]()

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.unsafeFind(key); ok {
		return fmt.Errorf("resource of type %v is already registered", key)
	}

	r.resourceList[key] = t
	return nil
}

func addAny(r *registry, res any) error {
	if res == nil {
		return fmt.Errorf("cannot add nil resource")
	}

	key := reflect.TypeOf(res)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.unsafeFind(key); ok {
		return fmt.Errorf("resource of type %v is already registered", key)
	}

	r.resourceList[key] = res
	return nil
}

func Get[T any]() (T, bool) {
	return get[T](globalRegistry)
}

func get[T any](r *registry) (T, bool) {
	var zero T
	key := keyOf[T]()

	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.unsafeFind(key)
	if !ok {
		return zero, false
	}

	return val.(T), true
}

func Walk(fn func(t reflect.Type, res any) bool) {
	globalRegistry.walk(fn)
}

type registry struct {
	mu           sync.RWMutex
	resourceList map[reflect.Type]any
}

func newRegistry() *registry {
	return &registry{
		resourceList: make(map[reflect.Type]any),
	}
}

func (r *registry) walk(fn func(t reflect.Type, res any) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for t, res := range r.resourceList {
		if !fn(t, res) {
			break
		}
	}
}

func (r *registry) unsafeFind(t reflect.Type) (any, bool) {
	val, ok := r.resourceList[t]
	return val, ok
}
