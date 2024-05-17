package in_memory

import "github.com/lazylex/watch-store/secure/internal/errors"

const inMemoryType = "in memory repo"

var (
	ErrNotNumericValue = NewInMemoryError("not numeric value")
)

// FullInMemoryError возвращает полностью заполненную структуру с типом InMemoryType.
func FullInMemoryError(message, origin string, initialError error) *errors.BaseError {
	return &errors.BaseError{
		Type:         inMemoryType,
		Message:      message,
		Origin:       origin,
		InitialError: initialError,
	}
}

// NewInMemoryError возвращает структуру ошибки с типом InMemoryType и переданным в качестве аргумента сообщением.
func NewInMemoryError(message string) *errors.BaseError {
	return &errors.BaseError{Type: inMemoryType, Message: message}
}
