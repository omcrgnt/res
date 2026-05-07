package res

import (
	"fmt"
	"testing"
)

// 1. Тестовые билдеры и ресурсы
type MyResource struct{ Value string }

type SystemBuilder struct{ Config string }

func (b *SystemBuilder) Build() (any, error) {
	return &MyResource{Value: "system: " + b.Config}, nil
}

type UserBuilder struct{ Config string }

func (b *UserBuilder) Build() (any, error) {
	return &MyResource{Value: "user: " + b.Config}, nil
}

// 2. Тест перекрытия (Override)
func TestFactory_Override(t *testing.T) {
	f := newFactory()

	// Регистрируем системный билдер
	f.register(&SystemBuilder{Config: "default"})

	// В пользовательском источнике передаем билдер того же типа, но с другим конфигом
	source := struct {
		MyBuilder Builder
	}{
		MyBuilder: &SystemBuilder{Config: "custom"},
	}

	reg, err := f.WithSource(source).build()
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	res, ok := get[*MyResource](reg)
	if !ok {
		t.Fatal("resource not found")
	}

	expected := "system: custom"
	if res.Value != expected {
		t.Errorf("expected %q, got %q", expected, res.Value)
	}
}

// 3. Тест на дубликаты РЕСУРСОВ (когда разные билдеры возвращают один тип)
type AnotherBuilder struct{}

func (b *AnotherBuilder) Build() (any, error) {
	return &MyResource{Value: "conflict"}, nil
}

func TestFactory_ResourceConflict(t *testing.T) {
	f := newFactory()

	// Два разных типа билдеров, но оба создают *MyResource
	f.register(&SystemBuilder{})

	source := struct {
		Conflict Builder
	}{
		Conflict: &AnotherBuilder{},
	}

	_, err := f.WithSource(source).build()
	if err == nil {
		t.Fatal("expected error due to resource type conflict, but got nil")
	}

	fmt.Println("Expected error caught:", err)
}
