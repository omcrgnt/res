package res

type MockRes struct{ Val string }
type MockBuilder struct{ V string }

func (b *MockBuilder) Build() (any, error) {
	return &MockRes{Val: b.V}, nil
}

// func TestBuild_Lifecycle(t *testing.T) {
// 	// Очищаем глобальное состояние перед тестом
// 	gf = newFactory()

// 	// 1. Системная регистрация
// 	gf.register(&MockBuilder{V: "system"})

// 	// 2. Пользовательская подмена (тот же тип билдера)
// 	cfg := struct {
// 		B *MockBuilder
// 	}{
// 		B: &MockBuilder{V: "user"},
// 	}

// 	// 3. Запуск сборки
// 	reg, err := gf.withSource(cfg).run()
// 	if err != nil {
// 		t.Fatalf("Build failed: %v", err)
// 	}

// 	// 4. Проверка результата в глобальном реестре
// 	res, ok := get[*MockRes](reg)
// 	if !ok {
// 		t.Fatal("Resource not found after Build")
// 	}

// 	if res.Val != "user" {
// 		t.Errorf("Expected 'user' (override), got %v", res.Val)
// 	}
// }
