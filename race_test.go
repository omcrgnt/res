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
		go func() {
			defer wg.Done()
			_ = Global().Add(stubNew{})
		}()
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = Global().AddWithTags(stubBuild{}, TagReplaceable)
		}()
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = Global().Transform(func(r any) any { return r })
		}()
	}

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = Global().GetByType(reflect.TypeFor[stubNew]())
			_ = Global().GetByInterface(reflect.TypeFor[Shaper]())
			Global().WalkEntries(func(Entry) bool { return true })
		}()
	}

	wg.Wait()
}
