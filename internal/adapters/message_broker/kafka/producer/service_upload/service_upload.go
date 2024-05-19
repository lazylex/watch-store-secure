package service_upload

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/errors/message_broker"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"time"
)

// ServiceUpload посылает сообщение о загрузке данного сервиса, что для других сервисов в системе должно означать, что
// их токены удалены из памяти и, следовательно, их следует запросить заново.
func ServiceUpload(cfg *config.Kafka) {
	var err error
	const retries = 3

	w := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  cfg.NeedToUpdateTokenTopic,
		AllowAutoTopicCreation: true,
	}

	defer func() {
		if err = w.Close(); err != nil {
			slog.Error(message_broker.ErrFailedToCloseWriter.Error())
		}
	}()

	origin := "ServiceUpload"

	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = w.WriteMessages(ctx, kafka.Message{Value: []byte("service upload")})
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			slog.Error(message_broker.FullMessageBrokerError("unexpected error", origin, err).Error())
			return
		}
		break
	}

	if err != nil {
		slog.Error(message_broker.ErrCouldNotSendMessage.WithOrigin(origin).Error())
	} else {
		slog.Info("kafka: successfully sent message about uploaded service")
	}
}
