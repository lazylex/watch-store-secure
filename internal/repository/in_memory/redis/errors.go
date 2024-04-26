package redis

import (
	"github.com/lazylex/watch-store/secure/internal/errors"
	"github.com/lazylex/watch-store/secure/internal/errors/in_memory"
	"strings"
)

const originPlace = "redis → "

// adaptErr переводит пришедшую ошибку к структурированной ошибке InMemory
func adaptErr(err error) error {
	return adaptErrSkipFrames(err, 2)
}

// adaptErrSkipFrames переводит пришедшую ошибку к структурированной ошибке InMemory с учетом последовательности
// вызова функций
func adaptErrSkipFrames(err error, skip int) error {
	if err == nil {
		return nil
	}
	origin := errors.GetFrame(skip).Function
	origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]

	return in_memory.FullInMemoryError("", origin, err)
}

// withOrigin добавляет место генерации ошибки
func withOrigin(err *errors.BaseError) error {
	origin := errors.GetFrame(2).Function
	origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]

	return err.WithOrigin(origin)
}

// ErrNotNumericValue возвращает ошибку in_memory.ErrNotNumericValue с местом генерации ошибки
func ErrNotNumericValue() error {
	return withOrigin(in_memory.ErrNotNumericValue)
}
