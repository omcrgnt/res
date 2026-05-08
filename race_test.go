package res

import (
	"sync"
	"testing"
)

type RaceBuilder struct{}

func (b *RaceBuilder) Build() (any, error) { return "ok", nil }

func TestRaceCondition(t *testing.T) {
	// Сбрасываем состояние
	gf = newFactory()
	globalRegistry = newRegistry()

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers * 4)

	// 1. Поток регистраторов
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			Register(&RaceBuilder{})
		}()
	}

	// 2. Поток сборщиков
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_ = Build(nil)
		}()
	}

	// 3. Поток читателей
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _ = Get[string]()
		}()
	}

	// 4. Поток файнедеров
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			// Ищем интерфейсы, пока всё перестраивается
			_ = Find[Shaper]()
		}()
	}

	wg.Wait()
}
