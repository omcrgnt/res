package res

import (
	"fmt"
	"reflect"
)

// TransformFunc преобразует ресурс перед wiring (например, obs-обёртка).
type TransformFunc func(any) any

// Transform применяет цепочку преобразований к ресурсам глобального реестра in-place.
func Transform(fns ...TransformFunc) error {
	return globalRegistry.transform(fns...)
}

func (r *registry) transform(fns ...TransformFunc) error {
	if len(fns) == 0 {
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for i, e := range r.resources {
		oldType := reflect.TypeOf(e.value)
		res := e.value

		for _, fn := range fns {
			res = fn(res)
		}

		newType := reflect.TypeOf(res)
		if newType != oldType {
			delete(r.byType, oldType)
			if _, exists := r.byType[newType]; exists {
				return fmt.Errorf("transform: type %v is already registered", newType)
			}
		}
		e.value = res
		r.byType[newType] = e
		r.resources[i] = e
	}

	return nil
}
