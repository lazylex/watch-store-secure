package service

import (
	"github.com/lazylex/watch-store/secure/internal/errors"
	"github.com/lazylex/watch-store/secure/internal/errors/joint"
	"github.com/lazylex/watch-store/secure/internal/errors/service"
	"strings"
)

const originPlace = "service → "

// adaptErr переводит пришедшую ошибку к структурированной ошибке.
func adaptErr(err error) error {
	return adaptErrSkipFrames(err, 2)
}

// adaptErrSkipFrames переводит пришедшую ошибку к структурированной ошибке с учетом последовательности
// вызова функций.
func adaptErrSkipFrames(err error, skip int) error {
	var message string

	if err == nil {
		return nil
	}

	origin := errors.GetFrame(skip).Function
	origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]

	if be, ok := err.(*errors.BaseError); ok {
		message = be.Message

		switch {
		case message == joint.ErrDuplicateData.Message:
			return service.ErrAlreadyExist.WithOrigin(be.Origin)
		case message == joint.ErrDataNotSaved.Message:
			return service.ErrNothingWasChanged.WithOrigin(be.Origin)
		case message == joint.ErrEmptyResult.Message:
			return service.ErrEmptyResult.WithOrigin(be.Origin)
		}

		if be.Type == service.ErrServiceType {
			return service.FullServiceError(be.Message, origin, be.InitialError)
		} else {
			return service.FullServiceError(be.Message, be.Origin, be.InitialError)
		}
	}

	return service.FullServiceError("", origin, err)
}

// withOrigin добавляет место генерации ошибки.
func withOrigin(err *errors.BaseError) error {
	origin := errors.GetFrame(2).Function
	origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]

	return err.WithOrigin(origin)
}

// ErrNotEnabledAccount возвращает ошибку service.ErrNotEnabledAccount с местом генерации ошибки.
func ErrNotEnabledAccount() error {
	return withOrigin(service.ErrNotEnabledAccount)
}

// ErrLogout возвращает ошибку service.ErrLogout с местом генерации ошибки.
func ErrLogout() error {
	return withOrigin(service.ErrLogout)
}
