package logger

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/config"
	"log"
	"log/slog"
	"os"
)

type ContextKey int

const (
	RequestId    ContextKey = 0
	TxId         ContextKey = 1
	OPLabel                 = "op"
	RequestLabel            = "request_id"
	TxLabel                 = "tx_number"
)

// MustCreate возвращает экземпляр *slog.Logger или останавливает программу, если окружение environment указано неверно
func MustCreate(environment, instance string) *slog.Logger {
	var logger *slog.Logger
	switch environment {
	case config.EnvironmentLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvironmentDebug:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		logger = logger.With(slog.String("instance", instance))
	case config.EnvironmentProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		logger = logger.With(slog.String("instance", instance))
	default:
		log.Fatal("program environment not set or it incorrect")
	}

	return logger
}

// LogWithCtxData извлекает, при наличии, из контекста идентификатор запроса и номер транзакции и добавляет в логгер
func LogWithCtxData(ctx context.Context, log *slog.Logger) *slog.Logger {
	if ctx.Value(RequestId) != nil {
		log = log.With(RequestLabel, ctx.Value(RequestId))
	}
	if ctx.Value(TxId) != nil {
		log = log.With(TxLabel, ctx.Value(TxId))
	}

	return log
}

// Null возвращает логгер, который пишет в никуда. Полезно для тестирования структур и функций, зависящих от логгера
func Null() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.NewFile(0, os.DevNull), nil))
}
