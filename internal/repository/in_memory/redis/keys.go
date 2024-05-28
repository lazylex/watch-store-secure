package redis

import (
	"fmt"
	"github.com/google/uuid"

	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
)

const (
	prefixSession                          = "s"
	prefixSessionByUUID                    = "si"
	prefixServicePermissionsNumbers        = "spn"
	prefixServicePermissionsNumbersForUser = "spn4u"
	prefixInstancePermissionsNumbers       = "ipn"
	prefixUuidHash                         = "uh"
	prefixAccountState                     = "as"
	prefixInstance                         = "i"
)

// keySession ключ для получения UUID пользователя сессии.
func keySession(sessionToken string) string {
	return fmt.Sprintf("%s:%s", prefixSession, sessionToken)
}

// keySessionByUUID ключ для получения токена сессии по UUID пользователя.
func keySessionByUUID(userID string) string {
	return fmt.Sprintf("%s:%s", prefixSessionByUUID, userID)
}

// keyServicePermissionsNumbersForUser ключ для получения списка разрешений сервиса service для пользователя (сервиса) с
// UUID равным id.
func keyServicePermissionsNumbersForUser(service string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", prefixServicePermissionsNumbersForUser, service, id.String())
}

// keyServicePermissionsNumbers ключ для получения нумерованного списка всех возможных разрешений сервиса.
func keyServicePermissionsNumbers(service string) string {
	return fmt.Sprintf("%s:%s", prefixServicePermissionsNumbers, service)
}

// keyInstancePermissionsNumbers ключ для получения списка разрешений экземпляра для пользователя (сервиса) с UUID
// равным id.
func keyInstancePermissionsNumbers(instance string, id uuid.UUID) string {
	return fmt.Sprintf("%s:%s:%s", prefixInstancePermissionsNumbers, instance, id.String())
}

// keyUserIdAndPasswordHash ключ для получения идентификатора пользователя и хэша его пароля по логину.
func keyUserIdAndPasswordHash(login loginVO.Login) string {
	return fmt.Sprintf("%s:%s", prefixUuidHash, string(login))
}

// keyAccountStateByLogin ключ для получения состояния учетной записи по логину.
func keyAccountStateByLogin(login loginVO.Login) string {
	return fmt.Sprintf("%s:%s", prefixAccountState, login)
}

// keyInstance ключ для получения секрета и названия сервиса для экземпляра.
func keyInstance(instance string) string {
	return fmt.Sprintf("%s:%s", prefixInstance, instance)
}
