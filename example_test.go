package res_test

import (
	"fmt"

	"github.com/omcrgnt/res"
)

// Logger — это ресурс, который мы будем использовать в приложении.
type Logger struct {
	Level string
}

// LoggerBuilder реализует интерфейс res.Builder.
// Он отвечает за создание ресурса Logger.
type LoggerBuilder struct {
	DefaultLevel string
}

func (b *LoggerBuilder) Build() (any, error) {
	return &Logger{Level: b.DefaultLevel}, nil
}

func Example() {
	// ВАЖНО: сбросим состояние, если тесты запускаются кучей
	// (в реальном приложении это не нужно, там Build вызывается один раз)

	res.Register(&LoggerBuilder{DefaultLevel: "INFO"})

	cfg := struct {
		Log *LoggerBuilder // Тип должен быть в точности *LoggerBuilder
	}{
		Log: &LoggerBuilder{DefaultLevel: "DEBUG"},
	}

	_ = res.Build(cfg)

	logger, _ := res.Get[*Logger]()
	fmt.Println(logger.Level)

	// Output: DEBUG
}

func ExampleGet() {
	// Простой пример добавления и получения ресурса напрямую.
	res.Add("my-secret-key")

	val, ok := res.Get[string]()
	if ok {
		fmt.Println("Found:", val)
	}
	// Output: Found: my-secret-key
}
