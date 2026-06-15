package res

import (
	"reflect"
	"sync"
	"testing"
)

func TestRaceCondition(t *testing.T) {
	resetGlobalRegistry()

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers * 4)

	for i := 0; i < workers; i++ {
		go func(n int) {
			defer wg.Done()
			_ = Add(n)
		}(i)
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = AddWithTags("bulk", TagReplaceable)
		}()
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = Transform(func(r any) any { return r })
		}()
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = GetByType(reflect.TypeFor[int]())
			_ = GetByInterface(reflect.TypeFor[Shaper]())
			WalkEntries(func(Entry) bool { return true })
		}()
	}

	wg.Wait()
}
