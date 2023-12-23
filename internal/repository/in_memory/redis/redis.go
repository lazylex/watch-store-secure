package redis

import (
	"github.com/go-redis/redis"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"time"
)

const (
	token = "token:"
)

type Redis struct {
	client *redis.Client
}

func Create(address, password string, db int) *Redis {
	client := redis.NewClient(&redis.Options{Addr: address, Password: password, DB: db})
	return &Redis{client: client}
}

// SaveSession сохраняет данные о времени жизни сессии и её пользователе. Переданный токен служит для создания ключа,
// по которому хранится идентификатор пользователя (хранение осуществляется переданное в TTL количество секунд)
func (r *Redis) SaveSession(dto dto.SessionDTO) error {
	_, err := r.client.Set(token+dto.Token, int(dto.Id), time.Duration(float64(dto.TTL))*time.Second).Result()
	return err
}
