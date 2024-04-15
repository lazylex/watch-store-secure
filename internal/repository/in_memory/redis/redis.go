package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/logger"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Redis struct {
	client *redis.Client
}

const (
	userIdField = "user_id"
	hashField   = "hash"
)

var (
	ErrNotNumericValue = redisErr("not numeric value")
	ErrIncorrectState  = redisErr("incorrect state")
)

// redisErr возвращает ошибку с префиксом redis
func redisErr(text string) error {
	return errors.New("redis: " + text)
}

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

func (r *Redis) GetUserIdAndPasswordHash(ctx context.Context, login loginVO.Login) (dto.UserIdWithPasswordHashDTO, error) {
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

func (r *Redis) SetUserIdAndPasswordHash(ctx context.Context, data dto.UserLoginAndIdWithPasswordHashDTO) {
	key := userIdAndPasswordHashKey(data.Login)
	r.client.HSet(ctx, key, userIdField, data.UserId.String(), hashField, data.Hash)
	// TODO определить ttl из конфигурации
	r.client.Expire(ctx, key, 1*time.Hour)
}

// GetAccountStateByLogin возвращает состояние учетной записи с переданным логином
func (r *Redis) GetAccountStateByLogin(ctx context.Context, login loginVO.Login) (account_state.State, error) {
	var numericVal int
	var err error
	var val string

	key := accountStateByLoginKey(login)

	if val, err = r.client.Get(ctx, key).Result(); err != nil {
		return 0, err
	}

	if numericVal, err = strconv.Atoi(val); err != nil {
		return 0, ErrNotNumericValue
	}

	// TODO определить ttl из конфигурации
	defer r.client.Expire(ctx, key, 24*time.Hour)

	return account_state.State(numericVal), err
}

// SetAccountState сохраняет состояние аккаунта с переданным логином
func (r *Redis) SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error {
	if !account_state.IsStateCorrect(stateDTO.State) {
		return ErrIncorrectState
	}
	// TODO считывать ttl из конфигурации
	return r.client.Set(ctx, accountStateByLoginKey(stateDTO.Login), int(stateDTO.State), 24*time.Hour).Err()
}

// SetServicePermissionsNumbersForAccount сохраняет номера разрешений аккаунта для сервиса
func (r *Redis) SetServicePermissionsNumbersForAccount(ctx context.Context, data dto.ServiceNameWithUserIdAndPermNumbersDTO) error {
	key := servicePermissionsNumbersKey(data.Service, data.UserId)
	return r.setPermissionsNumbers(ctx, key, data.PermissionNumbers)
}

// GetServicePermissionsNumbersForAccount возвращает номера всех разрешений аккаунта для сервиса
func (r *Redis) GetServicePermissionsNumbersForAccount(ctx context.Context, data dto.ServiceNameWithUserIdDTO) ([]int, error) {
	key := servicePermissionsNumbersKey(data.Service, data.UserId)
	return r.getPermissionsNumbers(ctx, key)
}

// SetInstancePermissionsNumbersForAccount сохраняет номера разрешений аккаунта для экземпляра сервиса
func (r *Redis) SetInstancePermissionsNumbersForAccount(ctx context.Context, data dto.ServiceNameWithUserIdAndPermNumbersDTO) error {
	key := instancePermissionsNumbersKey(data.Service, data.UserId)
	return r.setPermissionsNumbers(ctx, key, data.PermissionNumbers)
}

// GetInstancePermissionsNumbersForAccount возвращает номера разрешений аккаунта для экземпляра сервиса
func (r *Redis) GetInstancePermissionsNumbersForAccount(ctx context.Context, data dto.ServiceNameWithUserIdDTO) ([]int, error) {
	key := instancePermissionsNumbersKey(data.Service, data.UserId)
	return r.getPermissionsNumbers(ctx, key)
}

// setPermissionsNumbers сохраняет номера разрешений аккаунта по заданному ключу
func (r *Redis) setPermissionsNumbers(ctx context.Context, key string, permissionNumbers []int) error {

	numbers := make([]interface{}, len(permissionNumbers))

	for i, v := range permissionNumbers {
		numbers[i] = v
	}

	if err := r.client.SAdd(ctx, key, numbers...).Err(); err != nil {
		return err
	}

	// TODO считывать ttl из конфигурации
	return r.client.Expire(ctx, key, 24*time.Hour).Err()

}

// getPermissionsNumbers возвращает номера разрешений аккаунта по заданному ключу
func (r *Redis) getPermissionsNumbers(ctx context.Context, key string) ([]int, error) {
	if values, err := r.client.SMembers(ctx, key).Result(); err != nil {
		return nil, err
	} else {
		result := make([]int, 0, len(values))
		for _, v := range values {
			if permNumber, err := strconv.Atoi(v); err == nil {
				result = append(result, permNumber)
			}
		}

		// TODO считывать ttl из конфигурации
		return result, r.client.Expire(ctx, key, 24*time.Hour).Err()
	}
}

// ExistServicePermissionsNumbersForAccount возвращает true, если в памяти сохранены номера разрешений сервиса для
// аккаунта
func (r *Redis) ExistServicePermissionsNumbersForAccount(ctx context.Context, data dto.ServiceNameWithUserIdDTO) bool {
	key := servicePermissionsNumbersKey(data.Service, data.UserId)
	if result, err := r.client.Exists(ctx, key).Result(); err != nil {
		return false
	} else if result == 0 {
		return false
	}
	return true
}

// sessionKey ключ для получения UUID пользователя сессии
func sessionKey(sessionToken string) string {
	return fmt.Sprintf("session:%s", sessionToken)
}

// servicePermissionsNumbersKey ключ для получения списка разрешений сервиса service для пользователя (сервиса) с UUID равным id
func servicePermissionsNumbersKey(service string, id uuid.UUID) string {
	return fmt.Sprintf("serv_perm_numbs:%s:%s", service, id.String())
}

// instancePermissionsNumbersKey ключ для получения списка разрешений экземпляра для пользователя (сервиса) с UUID равным id
func instancePermissionsNumbersKey(instance string, id uuid.UUID) string {
	return fmt.Sprintf("inst_perm_numbs:%s:%s", instance, id.String())
}

// userIdAndPasswordHashKey ключ для получения идентификатора пользователя и хэша его пароля по логину
func userIdAndPasswordHashKey(login loginVO.Login) string {
	return fmt.Sprintf("uuid:hash:%s", string(login))
}

// accountStateByLoginKey ключ для получения состояния учетной записи по логину
func accountStateByLoginKey(login loginVO.Login) string {
	return fmt.Sprintf("account_state:%s", login)
}
