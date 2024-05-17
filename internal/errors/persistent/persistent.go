package persistent

import "github.com/lazylex/watch-store/secure/internal/errors"

const persistentType = "persistent repo"

var (
	ErrDuplicateKeyValue = NewPersistentError("duplicate key value violates unique constraint violation")
	ErrZeroRowsAffected  = NewPersistentError("zero rows affected")
	ErrNoRowsInResultSet = NewPersistentError("no rows in result set")
)

// FullPersistentError возвращает полностью заполненную структуру с типом PersistentType.
func FullPersistentError(message, origin string, initialError error) *errors.BaseError {
	return &errors.BaseError{
		Type:         persistentType,
		Message:      message,
		Origin:       origin,
		InitialError: initialError,
	}
}

// NewPersistentError возвращает структуру ошибки с типом PersistentType и переданным в качестве аргумента сообщением.
func NewPersistentError(message string) *errors.BaseError {
	return &errors.BaseError{Type: persistentType, Message: message}
}
