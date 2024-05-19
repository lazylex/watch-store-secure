package kafka

import (
	"github.com/lazylex/watch-store/secure/internal/adapters/message_broker/kafka/producer/service_upload"
	"github.com/lazylex/watch-store/secure/internal/config"
	"log/slog"
	"os"
)

// MustRun предназначен для запуска потребителей/продюсеров Кафки. Если в конфигурации cfg не задано имя топика, то
// соответствующий ему потребитель/продюсер не будет запущен. Работа приложения будет продолжена.
func MustRun(cfg *config.Kafka) {
	if len(cfg.Brokers) < 1 {
		slog.Error("empty kafka brokers list")
		os.Exit(1)
	}

	go service_upload.ServiceUpload(cfg)
}
