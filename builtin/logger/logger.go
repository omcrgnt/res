package logger

import (
	"context"
	"io"
	"log/slog"
	"os"

	logger "github.com/omcrgnt/proto/gen/go/logger/v1"
)

type Config struct {
	Level  logger.Level
	Format logger.Format
	Output io.Writer
}

func DefaultConfig() *Config {
	return &Config{
		Level:  logger.Level{Value: "debug"},
		Format: logger.Format{Value: "json"},
	}
}

func (c *Config) Build() (any, error) {
	if c.Output == nil {
		c.Output = os.Stdout
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(c.Level.Value)); err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if c.Format.Value == "json" {
		handler = slog.NewJSONHandler(c.Output, opts)
	} else {
		handler = slog.NewTextHandler(c.Output, opts)
	}

	// Возвращаем адаптер, чтобы соответствовать нашему интерфейсу Logger
	return &adapter{l: slog.New(handler)}, nil
}

// Адаптер для slog
type adapter struct {
	l *slog.Logger
}

func (a *adapter) Debug(ctx context.Context, msg string, args ...any) {
	a.l.DebugContext(ctx, msg, args...)
}
func (a *adapter) Info(ctx context.Context, msg string, args ...any) {
	a.l.InfoContext(ctx, msg, args...)
}
func (a *adapter) Warn(ctx context.Context, msg string, args ...any) {
	a.l.WarnContext(ctx, msg, args...)
}
func (a *adapter) Error(ctx context.Context, msg string, args ...any) {
	a.l.ErrorContext(ctx, msg, args...)
}

func (a *adapter) With(args ...any) Logger {
	return &adapter{l: a.l.With(args...)}
}
