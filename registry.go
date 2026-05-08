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

// Add добавляет готовый ресурс в глобальный реестр напрямую.
// Возвращает ошибку, если ресурс такого типа уже зарегистрирован.
func Add[T any](t T) error {
	gf.mu.Lock()
	defer gf.mu.Unlock()

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

// Get извлекает ресурс из глобального реестра по его типу.
// Если запрашивается интерфейс, функция попытается найти в реестре объект,
// реализующий этот интерфейс. Возвращает (resource, true), если поиск успешен.
func Get[T any]() (T, bool) {
	gf.mu.RLock()
	r := globalRegistry
	gf.mu.RUnlock()

	return get[T](r)
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

// Walk выполняет обход всех зарегистрированных ресурсов в глобальном реестре.
// Функция fn вызывается для каждого ресурса; если она возвращает false, обход прерывается.
func Walk(fn func(t reflect.Type, res any) bool) {
	gf.mu.RLock()
	r := globalRegistry
	gf.mu.RUnlock()

	r.walk(fn)
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

// Find возвращает список всех ресурсов, соответствующих типу T.
// Если T — интерфейс, функция вернет все ресурсы, реализующие его.
// Если T — не интерфейс (конкретный тип: структура, строка, число и т.д.),
// вернется срез, содержащий максимум один элемент.
func Find[T any]() []T {
	gf.mu.RLock()
	r := globalRegistry
	gf.mu.RUnlock()

	var matches []T
	targetType := reflect.TypeFor[T]()

	r.walk(func(t reflect.Type, res any) bool {
		switch targetType.Kind() {
		case reflect.Interface:
			// Если ищем интерфейс — проверяем, реализует ли его тип ресурса
			if t.Implements(targetType) {
				if val, ok := res.(T); ok {
					matches = append(matches, val)
				}
			}
		default:
			// Для всех остальных типов (структуры, строки и т.д.) — только точное совпадение
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
