package redis

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"os"
	"strconv"
)

type Redis struct {
	client *redis.Client
	ttl    config.TTL
}

const (
	userIdField = "user_id"
	hashField   = "hash"
)

// MustCreate создание структуры с клиентом для взаимодействия с Redis. При ошибке соединения с сервером Redis выводит
// ошибку в лог и прекращает работу приложения
func MustCreate(cfg config.Redis, ttl config.TTL) *Redis {
	client := redis.NewClient(
		&redis.Options{Addr: cfg.RedisAddress, Username: cfg.RedisUser, Password: cfg.RedisPassword, DB: cfg.RedisDB})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		slog.Error(adaptErr(err).Error())
		os.Exit(1)
	} else {
		slog.Info("successfully received pong from redis server")
	}

	return &Redis{client: client, ttl: ttl}
}

// SaveSession сохраняет данные о времени жизни сессии и её пользователе. Переданный токен служит для создания ключа,
// по которому хранится идентификатор пользователя (хранение осуществляется переданное в TTL количество секунд)
func (r *Redis) SaveSession(ctx context.Context, dto *dto.UserIdToken) error {
	pipe := r.client.Pipeline()
	pipe.Set(ctx, keySession(dto.Token), dto.UserId.String(), r.ttl.SessionTTL)
	pipe.Set(ctx, keySessionByUUID(dto.UserId.String()), dto.Token, r.ttl.SessionTTL)
	_, err := pipe.Exec(ctx)

	return adaptErr(err)
}

// DeleteSession удаляет из памяти данные о привязке токена к UUID пользователя и привязке UUID пользователя к токену
// сессии
func (r *Redis) DeleteSession(ctx context.Context, id uuid.UUID) error {
	sessionByUUID := keySessionByUUID(id.String())
	sessionToken, err := r.client.Get(ctx, sessionByUUID).Result()
	if err != nil {
		return adaptErr(err)
	}
	session := keySession(sessionToken)

	return adaptErr(r.client.Del(ctx, sessionByUUID, session).Err())
}

// IsSessionActiveByUUID возвращает true, если существует сессия для пользователя (сервиса) с переданным идентификатором
func (r *Redis) IsSessionActiveByUUID(ctx context.Context, userId uuid.UUID) bool {
	key := keySessionByUUID(userId.String())
	if result, err := r.client.Exists(ctx, key).Result(); err != nil {
		return false
	} else {
		return result == 1
	}
}

// IsSessionActiveByToken возвращает true, если существует сессия для пользователя (сервиса) с переданным токеном сессии
func (r *Redis) IsSessionActiveByToken(ctx context.Context, token string) bool {
	key := keySession(token)
	if result, err := r.client.Exists(ctx, key).Result(); err != nil {
		return false
	} else {
		return result == 1
	}
}

// extendSessionLife продлевает жизнь данным по ключам, относящимся к сессии пользователя (сервиса)
func (r *Redis) extendSessionLife(ctx context.Context, key string) error {
	var mUUID string
	var err error

	if mUUID, err = r.client.Get(ctx, key).Result(); err != nil {
		return adaptErr(err)
	}

	pipe := r.client.Pipeline()
	pipe.Expire(ctx, key, r.ttl.SessionTTL)
	pipe.Expire(ctx, keySessionByUUID(mUUID), r.ttl.SessionTTL)
	_, err = pipe.Exec(ctx)

	return adaptErr(err)
}

// GetUserUUIDFromSession получает UUID пользователя сессии
func (r *Redis) GetUserUUIDFromSession(ctx context.Context, sessionToken string) (uuid.UUID, error) {
	var val []byte
	var err error
	var parsedUUID uuid.UUID
	key := keySession(sessionToken)
	if val, err = r.client.Get(ctx, key).Bytes(); err != nil {
		return uuid.Nil, adaptErr(err)
	}

	defer func(r *Redis, ctx context.Context, key string) {
		_ = r.extendSessionLife(ctx, key)
	}(r, ctx, key)

	parsedUUID, err = uuid.FromBytes(val)

	return parsedUUID, adaptErr(err)
}

func (r *Redis) GetUserIdAndPasswordHash(ctx context.Context, login loginVO.Login) (dto.UserIdHash, error) {
	var err error
	var parsedUUID uuid.UUID
	var values map[string]string

	key := keyUserIdAndPasswordHash(login)

	if values, err = r.client.HGetAll(ctx, key).Result(); err != nil {
		return dto.UserIdHash{}, adaptErr(err)
	}

	r.client.Expire(ctx, key, r.ttl.UserIdAndPasswordHashTTL)
	if parsedUUID, err = uuid.Parse(values[userIdField]); err != nil {
		return dto.UserIdHash{}, adaptErr(err)
	}
	return dto.UserIdHash{UserId: parsedUUID, Hash: values[hashField]}, nil
}

func (r *Redis) SetUserIdAndPasswordHash(ctx context.Context, data *dto.UserIdLoginHash) {
	key := keyUserIdAndPasswordHash(data.Login)
	r.client.HSet(ctx, key, userIdField, data.UserId.String(), hashField, data.Hash)
	r.client.Expire(ctx, key, r.ttl.UserIdAndPasswordHashTTL)
}

// GetAccountStateByLogin возвращает состояние учетной записи с переданным логином
func (r *Redis) GetAccountStateByLogin(ctx context.Context, login loginVO.Login) (account_state.State, error) {
	var numericVal int
	var err error
	var val string

	key := keyAccountStateByLogin(login)

	if val, err = r.client.Get(ctx, key).Result(); err != nil {
		return 0, adaptErr(err)
	}

	if numericVal, err = strconv.Atoi(val); err != nil {
		return 0, ErrNotNumericValue()
	}

	defer r.client.Expire(ctx, key, r.ttl.AccountStateTTL)

	return account_state.State(numericVal), nil
}

// SetAccountState сохраняет состояние аккаунта с переданным логином
func (r *Redis) SetAccountState(ctx context.Context, data *dto.LoginState) error {
	return adaptErr(r.client.Set(ctx, keyAccountStateByLogin(data.Login), int(data.State), r.ttl.AccountStateTTL).Err())
}

// SetServicePermissionsNumbersForAccount сохраняет номера разрешений аккаунта для сервиса
func (r *Redis) SetServicePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdServicePermNumbers) error {
	key := keyServicePermissionsNumbers(data.Service, data.UserId)
	return adaptErr(r.setPermissionsNumbers(ctx, key, data.PermissionNumbers))
}

// GetServicePermissionsNumbersForAccount возвращает номера всех разрешений аккаунта для сервиса
func (r *Redis) GetServicePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdService) ([]int, error) {
	key := keyServicePermissionsNumbers(data.Service, data.UserId)
	if numbers, err := r.getPermissionsNumbers(ctx, key); err != nil {
		return []int{}, adaptErr(err)
	} else {
		return numbers, nil
	}

}

// SetInstancePermissionsNumbersForAccount сохраняет номера разрешений аккаунта для экземпляра сервиса
func (r *Redis) SetInstancePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdInstancePermNumbers) error {
	key := keyInstancePermissionsNumbers(data.Instance, data.UserId)
	return adaptErr(r.setPermissionsNumbers(ctx, key, data.PermissionNumbers))
}

// GetInstancePermissionsNumbersForAccount возвращает номера разрешений аккаунта для экземпляра сервиса
func (r *Redis) GetInstancePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdInstance) ([]int, error) {
	key := keyInstancePermissionsNumbers(data.Instance, data.UserId)
	if numbers, err := r.getPermissionsNumbers(ctx, key); err != nil {
		return []int{}, adaptErr(err)
	} else {
		return numbers, nil
	}
}

// SetInstanceServiceAndSecret сохраняет название сервиса для экземпляра и секретный ключ для создания подписи
// JWT-токена
func (r *Redis) SetInstanceServiceAndSecret(ctx context.Context, data *dto.NameServiceSecret) error {
	key := keyInstance(data.Name)
	if err := r.client.HSet(ctx, key, "service", data.Service, "secret", data.Secret).Err(); err != nil {
		return adaptErr(err)
	}

	defer r.client.Expire(ctx, key, r.ttl.InstanceDataTTL)

	return nil
}

// GetServiceName возвращает название сервиса по имени его экземпляра
func (r *Redis) GetServiceName(ctx context.Context, instanceName string) (string, error) {
	key := keyInstance(instanceName)
	result, err := r.client.HGet(ctx, key, "service").Result()
	if err != nil {
		return "", adaptErr(err)
	}

	defer r.client.Expire(ctx, key, r.ttl.InstanceDataTTL)

	return result, nil
}

// SetInstanceServiceName сохраняет название сервиса для экземпляра
func (r *Redis) SetInstanceServiceName(ctx context.Context, data *dto.NameService) error {
	key := keyInstance(data.Name)
	if err := r.client.HSet(ctx, key, "service", data.Service).Err(); err != nil {
		return adaptErr(err)
	}

	defer r.client.Expire(ctx, key, r.ttl.InstanceDataTTL)

	return nil
}

// GetInstanceSecret возвращает секретный ключ для экземпляра сервиса
func (r *Redis) GetInstanceSecret(ctx context.Context, instanceName string) (string, error) {
	key := keyInstance(instanceName)
	result, err := r.client.HGet(ctx, key, "secret").Result()
	if err != nil {
		return "", adaptErr(err)
	}

	defer r.client.Expire(ctx, key, r.ttl.InstanceDataTTL)

	return result, nil
}

// SetInstanceSecret сохраняет секретный ключ экземпляра сервиса, необходимый для создания подписи JWT-токена
func (r *Redis) SetInstanceSecret(ctx context.Context, data *dto.NameSecret) error {
	key := keyInstance(data.Name)
	if err := r.client.HSet(ctx, key, "secret", data.Secret).Err(); err != nil {
		return adaptErr(err)
	}

	defer r.client.Expire(ctx, key, r.ttl.InstanceDataTTL)

	return nil
}

// setPermissionsNumbers сохраняет номера разрешений аккаунта по заданному ключу
func (r *Redis) setPermissionsNumbers(ctx context.Context, key string, permissionNumbers []int) error {

	numbers := make([]interface{}, len(permissionNumbers))

	for i, v := range permissionNumbers {
		numbers[i] = v
	}

	if err := r.client.SAdd(ctx, key, numbers...).Err(); err != nil {
		return adaptErr(err)
	}

	return adaptErr(r.client.Expire(ctx, key, r.ttl.PermissionsNumbersTTL).Err())
}

// getPermissionsNumbers возвращает номера разрешений аккаунта по заданному ключу
func (r *Redis) getPermissionsNumbers(ctx context.Context, key string) ([]int, error) {
	if values, err := r.client.SMembers(ctx, key).Result(); err != nil {
		return []int{}, adaptErr(err)
	} else {
		result := make([]int, 0, len(values))
		for _, v := range values {
			if permNumber, err := strconv.Atoi(v); err == nil {
				result = append(result, permNumber)
			}
		}

		return result, adaptErr(r.client.Expire(ctx, key, r.ttl.PermissionsNumbersTTL).Err())
	}
}

// ExistServicePermissionsNumbersForAccount возвращает true, если в памяти сохранены номера разрешений сервиса для
// аккаунта
func (r *Redis) ExistServicePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdService) bool {
	key := keyServicePermissionsNumbers(data.Service, data.UserId)
	return r.existKey(ctx, key)
}

// ExistInstancePermissionsNumbersForAccount возвращает true, если в памяти сохранены номера разрешений экземпляра для
// аккаунта
func (r *Redis) ExistInstancePermissionsNumbersForAccount(ctx context.Context, data *dto.UserIdInstance) bool {
	key := keyInstancePermissionsNumbers(data.Instance, data.UserId)
	return r.existKey(ctx, key)
}

// existKey возвращает true, если по переданному ключу в памяти есть данные
func (r *Redis) existKey(ctx context.Context, key string) bool {
	if result, err := r.client.Exists(ctx, key).Result(); err != nil {
		return false
	} else if result == 0 {
		return false
	}

	return true
}
