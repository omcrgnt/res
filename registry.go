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
	return addAny(globalRegistry, v)
}

func (globalAccessor) Walk(fn func(t reflect.Type, res any) bool) {
	globalRegistry.walk(fn)
}

func keyOf[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

// Add добавляет готовый ресурс в глобальный реестр.
// Возвращает ошибку, если ресурс такого типа уже зарегистрирован.
func Add[T any](t T) error {
	return add(globalRegistry, t)
}

// AddAll добавляет несколько ресурсов в глобальный реестр.
func AddAll(resources ...any) error {
	for _, res := range resources {
		if err := addAny(globalRegistry, res); err != nil {
			return err
		}
	}
	return nil
}

func add[T any](r *registry, t T) error {
	return addAny(r, t)
}

func addAny(r *registry, res any) error {
	if res == nil {
		return fmt.Errorf("cannot add nil resource")
	}

	key := reflect.TypeOf(res)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byType[key]; ok {
		return fmt.Errorf("resource of type %v is already registered", key)
	}

	r.byType[key] = res
	r.resources = append(r.resources, res)
	return nil
}

// Walk выполняет read-only обход ресурсов в порядке регистрации.
func Walk(fn func(t reflect.Type, res any) bool) {
	globalRegistry.walk(fn)
}

// Get извлекает ресурс из глобального реестра по его типу.
func Get[T any]() (T, bool) {
	return get[T](globalRegistry)
}

func get[T any](r *registry) (T, bool) {
	var zero T
	key := keyOf[T]()

	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.byType[key]
	if !ok {
		return zero, false
	}

	return val.(T), true
}

type registry struct {
	mu        sync.RWMutex
	resources []any
	byType    map[reflect.Type]any
}

func newRegistry() *registry {
	return &registry{
		byType: make(map[reflect.Type]any),
	}
}

func (r *registry) walk(fn func(t reflect.Type, res any) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, res := range r.resources {
		if !fn(reflect.TypeOf(res), res) {
			break
		}
	}
}

// Find возвращает ресурсы, соответствующие типу T.
// Для интерфейса — все реализации; для concrete type — максимум один элемент.
func Find[T any]() []T {
	var matches []T
	targetType := reflect.TypeFor[T]()

	globalRegistry.walk(func(t reflect.Type, res any) bool {
		switch targetType.Kind() {
		case reflect.Interface:
			if t.Implements(targetType) {
				if val, ok := res.(T); ok {
					matches = append(matches, val)
				}
			}
		default:
			if t == targetType {
				if val, ok := res.(T); ok {
					matches = append(matches, val)
				}
			}
		}
		return true
	})

	return matches
}

func resetGlobalRegistry() {
	globalRegistry = newRegistry()
}
