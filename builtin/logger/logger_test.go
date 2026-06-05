package logger

import (
	"context"
	"testing"

	logger "github.com/omcrgnt/proto/gen/go/logger/v1"
)

func TestLogger_Build(t *testing.T) {
	cfg := &Config{
		Level:  logger.Level{Value: "debug"},
		Format: logger.Format{Value: "text"},
	}

	res, err := cfg.Build()
	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}

	l, ok := res.(Logger) // Проверяем соответствие интерфейсу
	if !ok {
		t.Fatal("result does not implement Logger interface")
	}

	// Проверяем, что вызов не падает
	l.Info(context.Background(), "test message", "key", "value")
}
