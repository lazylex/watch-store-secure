package in_memory

import "github.com/lazylex/watch-store/secure/internal/errors"

const inMemoryType = "in memory repo"

var (
	ErrNotNumericValue = NewInMemoryError("not numeric value")
	ErrIncorrectState  = NewInMemoryError("incorrect state")
)

type InMemory struct {
	errors.BaseError
}

// FullInMemoryError возвращает полностью заполненную структуру InMemory
func FullInMemoryError(message, origin string, initialError error) *InMemory {
	m := &InMemory{}
	m.Type = inMemoryType
	m.Message = message
	m.Origin = origin
	m.InitialError = initialError

	return m
}

// NewInMemoryError возвращает структуру ошибки InMemory с переданным в качестве аргумента сообщением
func NewInMemoryError(message string) *InMemory {
	m := &InMemory{}
	m.Type = inMemoryType
	m.Message = message

	return m
}

// WithOrigin добавляет в структуру место появления ошибки
func (m *InMemory) WithOrigin(origin string) *InMemory {
	m.Origin = origin
	return m
}
