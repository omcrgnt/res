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
			_ = AddAll("bulk")
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
			_, _ = Get[int]()
			_ = Find[Shaper]()
			Walk(func(reflect.Type, any) bool { return true })
		}()
	}

	wg.Wait()
}
