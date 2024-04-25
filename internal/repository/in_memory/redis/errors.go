package redis

import (
	"github.com/lazylex/watch-store/secure/internal/errors"
	"github.com/lazylex/watch-store/secure/internal/errors/in_memory"
	"strings"
)

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
	origin = "redis → " + origin[strings.LastIndex(origin, ".")+1:]

	return in_memory.FullInMemoryError("", origin, err)
}

func withOrigin(err *in_memory.InMemory) error {
	origin := errors.GetFrame(1).Function
	origin = "redis → " + origin[strings.LastIndex(origin, ".")+1:]

	return err.WithOrigin(origin)
}