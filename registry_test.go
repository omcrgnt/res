package res

import (
	"reflect"
	"testing"
)

func TestRegistry_AddGet(t *testing.T) {
	r := newRegistry()
	val := "hello"

	// Тестируем addAny (внутренний) через публичный интерфейс (если бы он был)
	// или просто напрямую, так как мы в пакете res
	err := addAny(r, val)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Тестируем get
	res, ok := get[string](r)
	if !ok || res != "hello" {
		t.Errorf("Get failed: expected 'hello', got %v", res)
	}

	// Тестируем дубликат
	err = addAny(r, "world")
	if err == nil {
		t.Error("Expected error on duplicate type registration, got nil")
	}
}

func TestRegistry_Walk(t *testing.T) {
	r := newRegistry()
	_ = addAny(r, 10)
	_ = addAny(r, "string")

	count := 0
	r.walk(func(t reflect.Type, res any) bool {
		count++
		return true
	})

	if count != 2 {
		t.Errorf("Walk failed: expected 2 items, got %d", count)
	}
}

// Определим интерфейс и реализацию для теста
type Shaper interface {
	Area() int
}

type Square struct {
	Side int
}

func (s *Square) Area() int {
	return s.Side * s.Side
}

func TestFind(t *testing.T) {
	// 1. Очищаем состояние
	gf = newFactory()
	globalRegistry = newRegistry()

	// 2. Добавляем ресурс напрямую (или через Build)
	sq := &Square{Side: 10}
	_ = Add(sq)

	// 3. ТЕСТ 1: Поиск по конкретному типу (должен найти 1 элемент)
	squares := Find[*Square]()
	if len(squares) != 1 {
		t.Errorf("Find[*Square] failed: expected 1 match, got %d", len(squares))
	}

	// 4. ТЕСТ 2: Поиск по интерфейсу (самое важное для SDI!)
	// Если твоя текущая реализация не использует Implements, этот тест УПАДЕТ.
	shapes := Find[Shaper]()
	if len(shapes) != 1 {
		t.Errorf("Find[Shaper] failed: expected 1 match (Square implements Shaper), got %d", len(shapes))
	}

	// 5. ТЕСТ 3: Поиск несуществующего типа
	strings := Find[string]()
	if len(strings) != 0 {
		t.Errorf("Find[string] failed: expected 0 matches, got %d", len(strings))
	}
}
