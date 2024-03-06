package redis

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/logger"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"time"
)

type Redis struct {
	client *redis.Client
}

// MustCreate создание структуры с клиентом для взаимодействия с Redis. При ошибке соеднинения с сервером Redis выводит
// ошибку в лог и прекращает работу приложения
func MustCreate(cfg config.Redis, log *slog.Logger) *Redis {
	log = log.With(logger.OPLabel, "repository.in_memory.redis.MustCreate")
	client := redis.NewClient(
		&redis.Options{Addr: cfg.RedisAddress, Username: cfg.RedisUser, Password: cfg.RedisPassword, DB: cfg.RedisDB})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	} else {
		log.Info("successfully received pong from redis server")
	}

	return &Redis{client: client}
}

// SaveSession сохраняет данные о времени жизни сессии и её пользователе. Переданный токен служит для создания ключа,
// по которому хранится идентификатор пользователя (хранение осуществляется переданное в TTL количество секунд)
func (r *Redis) SaveSession(ctx context.Context, dto dto.SessionDTO) error {
	_, err := r.client.Set(ctx, sessionKey(dto.Token), dto.UserId, time.Duration(float64(dto.TTL))*time.Second).Result()
	return err
}

// GetUserUUIDFromSession получает UUID пользователя сессии
func (r *Redis) GetUserUUIDFromSession(ctx context.Context, sessionToken string) (uuid.UUID, error) {
	var val []byte
	var err error

	if val, err = r.client.Get(ctx, sessionKey(sessionToken)).Bytes(); err != nil {
		return uuid.Nil, err
	}

	return uuid.FromBytes(val)
}

// sessionKey ключ для получения UUID пользователя сессии
func sessionKey(sessionToken string) string {
	return fmt.Sprintf("session:%s", sessionToken)
}

// permissionsKey ключ для получения списка разрешений сервиса service для пользователя (сервиса) с UUID равным id
func permissionsKey(service string, id uuid.UUID) string {
	return fmt.Sprintf("perm:%s:%s", service, id.String())
}
