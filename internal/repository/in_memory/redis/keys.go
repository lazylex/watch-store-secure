package redis

import (
	"fmt"
	"github.com/google/uuid"

	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
)

const (
	prefixSession                    = "s"
	prefixServicePermissionsNumbers  = "spn"
	prefixInstancePermissionsNumbers = "ipn"
	prefixUuidHash                   = "uh"
	prefixAccountState               = "as"
)

// keySession ключ для получения UUID пользователя сессии
func keySession(sessionToken string) string {
	return fmt.Sprintf("%s:%s", prefixSession, sessionToken)
}

// keyServicePermissionsNumbers ключ для получения списка разрешений сервиса service для пользователя (сервиса) с UUID равным id
func keyServicePermissionsNumbers(service string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", prefixServicePermissionsNumbers, service, id.String())
}

// keyInstancePermissionsNumbers ключ для получения списка разрешений экземпляра для пользователя (сервиса) с UUID равным id
func keyInstancePermissionsNumbers(instance string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", prefixInstancePermissionsNumbers, instance, id.String())
}

// keyUserIdAndPasswordHash ключ для получения идентификатора пользователя и хэша его пароля по логину
func keyUserIdAndPasswordHash(login loginVO.Login) string {
	return fmt.Sprintf("%s:%s", prefixUuidHash, string(login))
}

// keyAccountStateByLogin ключ для получения состояния учетной записи по логину
func keyAccountStateByLogin(login loginVO.Login) string {
	return fmt.Sprintf("%s:%s", prefixAccountState, login)
}