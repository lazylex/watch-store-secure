package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"time"
)

type Redis struct {
	client *redis.Client
}

// Create создание структуры с клиентом для взаимодействия с Redis
func Create(cfg config.Redis) *Redis {
	client := redis.NewClient(
		&redis.Options{Addr: cfg.RedisAddress, Password: cfg.RedisPassword, DB: cfg.RedisDB})
	return &Redis{client: client}
}

// SaveSession сохраняет данные о времени жизни сессии и её пользователе. Переданный токен служит для создания ключа,
// по которому хранится идентификатор пользователя (хранение осуществляется переданное в TTL количество секунд)
func (r *Redis) SaveSession(dto dto.SessionDTO) error {
	_, err := r.client.Set(sessionKey(dto.Token), dto.Id, time.Duration(float64(dto.TTL))*time.Second).Result()
	return err
}

// GetUserUUIDFromSession получает UUID пользователя сессии
func (r *Redis) GetUserUUIDFromSession(sessionToken string) (uuid.UUID, error) {
	var val []byte
	var err error

	if val, err = r.client.Get(sessionKey(sessionToken)).Bytes(); err != nil {
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
