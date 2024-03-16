package redis

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
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

const (
	userIdField = "user_id"
	hashField   = "hash"
)

// MustCreate создание структуры с клиентом для взаимодействия с Redis. При ошибке соеднинения с сервером Redis выводит
// ошибку в лог и прекращает работу приложения
func MustCreate(cfg config.Redis) *Redis {
	log := slog.With(logger.OPLabel, "repository.in_memory.redis.MustCreate")
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

func (r *Redis) GetUserIdAndPasswordHash(ctx context.Context, login login.Login) (dto.UserIdWithPasswordHashDTO, error) {
	var err error
	var parsedUUID uuid.UUID
	var values map[string]string

	key := userIdAndPasswordHashKey(login)

	if values, err = r.client.HGetAll(ctx, key).Result(); err != nil {
		return dto.UserIdWithPasswordHashDTO{}, err
	}

	// TODO определить ttl из конфигурации
	r.client.Expire(ctx, key, 1*time.Hour)
	if parsedUUID, err = uuid.Parse(values[userIdField]); err != nil {
		return dto.UserIdWithPasswordHashDTO{}, err
	}
	return dto.UserIdWithPasswordHashDTO{UserId: parsedUUID, Hash: values[hashField]}, nil
}

func (r *Redis) SetUserIdAndPasswordHash(ctx context.Context, data dto.UserLoginAndIdWithPasswordHashDTO) error {
	var err error
	key := userIdAndPasswordHashKey(data.Login)
	if err = r.client.HSet(ctx, key, userIdField, data.UserId.String(), hashField, data.Hash).Err(); err != nil {
		return err
	}
	// TODO определить ttl из конфигурации
	return r.client.Expire(ctx, key, 1*time.Hour).Err()
}

// sessionKey ключ для получения UUID пользователя сессии
func sessionKey(sessionToken string) string {
	return fmt.Sprintf("session:%s", sessionToken)
}

// permissionsKey ключ для получения списка разрешений сервиса service для пользователя (сервиса) с UUID равным id
func permissionsKey(service string, id uuid.UUID) string {
	return fmt.Sprintf("perm:%s:%s", service, id.String())
}

// userIdAndPasswordHashKey ключ для получения идентификатора пользователя и хэша его пароля
func userIdAndPasswordHashKey(login login.Login) string {
	return fmt.Sprintf("uuid:hash:%s", login)
}
