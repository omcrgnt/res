package res

import (
	"fmt"
	"reflect"
	"sync"
)

var globalRegistry = &Registry{
	resourceList: make(map[reflect.Type]any),
}

func keyOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func Add[T any](t T) error {
	key := keyOf[T]()

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if _, ok := globalRegistry.find(key); ok {
		return fmt.Errorf("resource of type %v is already registered", key)
	}

	globalRegistry.resourceList[key] = t
	return nil
}

func Get[T any]() (T, bool) {
	var zero T
	key := keyOf[T]()

	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	val, ok := globalRegistry.find(key)
	if !ok {
		return zero, false
	}

	return val.(T), true
}

type Registry struct {
	mu           sync.RWMutex
	resourceList map[reflect.Type]any
}

func (r *Registry) find(t reflect.Type) (any, bool) {
	val, ok := r.resourceList[t]
	return val, ok
}
