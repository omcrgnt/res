package res

import (
	"reflect"
	"testing"
)

// Определяем интерфейс и структуру вне функций теста
type Storage interface {
	Save() string
}

type mockStorage struct{}

func (m *mockStorage) Save() string { return "saved" }

// 1. Тест на изоляцию реестров (White Box)
func TestRegistry_Isolation(t *testing.T) {
	r1 := newRegistry()
	r2 := newRegistry()

	_ = add(r1, "resource-1")
	_ = add(r2, 42)

	if _, ok := get[string](r1); !ok {
		t.Error("r1 должен содержать строку")
	}
	if _, ok := get[int](r1); ok {
		t.Error("r1 НЕ должен содержать число")
	}
	if _, ok := get[int](r2); !ok {
		t.Error("r2 должен содержать число")
	}
}

// 2. Тест на работу с интерфейсами (самый важный кейс)
func TestRegistry_Interfaces(t *testing.T) {
	r := newRegistry()
	impl := &mockStorage{}

	// Регистрируем именно как интерфейс Storage
	err := add[Storage](r, impl)
	if err != nil {
		t.Fatalf("Ошибка регистрации интерфейса: %v", err)
	}

	// Пытаемся достать как интерфейс
	val, ok := get[Storage](r)
	if !ok {
		t.Fatal("Ресурс должен быть доступен по типу интерфейса")
	}

	if val.Save() != "saved" {
		t.Errorf("Ожидалось 'saved', получено %s", val.Save())
	}
}

// 3. Тест на ошибку при дубликате
func TestRegistry_Duplicate(t *testing.T) {
	r := newRegistry()

	_ = add(r, "first")
	err := add(r, "second") // тип тот же (string)

	if err == nil {
		t.Error("Ожидалась ошибка при повторной регистрации типа string")
	}
}

// 4. Тест публичного API (Global)
func TestGlobalAPI(t *testing.T) {
	type GlobalRes struct{ ID int }
	res := GlobalRes{ID: 777}

	err := Add(res)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	found, ok := Get[GlobalRes]()
	if !ok || found.ID != 777 {
		t.Errorf("Get вернул неверные данные: %+v", found)
	}
}

func TestRegistry_AnyTypeChallenge(t *testing.T) {
	r := newRegistry()

	// Имитируем выдачу билдера (он возвращает any)
	var rawResource any = "я реальная строка, а не просто any"

	// Сейчас мы вынуждены использовать add, но так как это generic,
	// если мы напишем add(r, rawResource), то T станет типом 'any'.
	err := addAny(r, rawResource)
	if err != nil {
		t.Fatalf("Не удалось добавить ресурс: %v", err)
	}

	// Пытаемся достать ресурс как строку.
	// ОЖИДАНИЕ: Мы должны найти строку.
	// РЕАЛЬНОСТЬ (сейчас): Мы не найдем её, так как в мапе ключ — интерфейс 'any'.
	val, ok := get[string](r)
	if !ok {
		t.Errorf("Провал! Ресурс был зарегистрирован как any, и мы не смогли достать его как string")
	}

	if val != "я реальная строка, а не просто any" {
		t.Errorf("Ожидалось содержимое строки, получено: %v", val)
	}
}

func TestRegistry_Walk(t *testing.T) {
	r := newRegistry()
	_ = add(r, "string-resource")
	_ = add(r, 100)
	_ = add(r, true)

	t.Run("Full scan", func(t *testing.T) {
		count := 0
		r.walk(func(typ reflect.Type, res any) bool {
			count++
			return true // продолжаем до конца
		})
		if count != 3 {
			t.Errorf("Ожидалось 3 ресурса, найдено %d", count)
		}
	})

	t.Run("Early break", func(t *testing.T) {
		count := 0
		r.walk(func(typ reflect.Type, res any) bool {
			count++
			return false // останавливаемся сразу после первого
		})
		if count != 1 {
			t.Errorf("Ожидалось прерывание после 1 ресурса, пройдено %d", count)
		}
	})
}
