package logger

import "context"

// Logger требует контекст для каждой операции.
// Это гарантирует, что TraceID и SpanID всегда будут в логах.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)

	// With возвращает новый Logger с добавленными полями.
	With(args ...any) Logger
}
