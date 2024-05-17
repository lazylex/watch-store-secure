/*
Package logger: пакет для создания логгера, соответствующего окружению, в котором запущено приложение. Функция
MustCreate возвращает указатель на slog.Logger. В случае несоответствия переданного в функцию окружения одному из
доступных вариантов (config.EnvironmentLocal, config.EnvironmentDebug или config.EnvironmentProduction) выполнение
приложения прекращается.
*/
package logger

import (
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/pkg/colorlog"
	"log"
	"log/slog"
	"os"
	"time"
)

// MustCreate возвращает экземпляр *slog.Logger или останавливает программу, если окружение environment указано неверно
func MustCreate(environment, instance string) *slog.Logger {
	var logger *slog.Logger
	switch environment {
	case config.EnvironmentLocal:
		logger = slog.New(colorlog.NewHandler(os.Stdout, &colorlog.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		}))
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
