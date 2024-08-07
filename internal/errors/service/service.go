package service

import "github.com/lazylex/watch-store/secure/internal/errors"

const ErrServiceType = "service"

var (
	ErrAuthenticationData = NewServiceError("incorrect login or password")
	ErrNotEnabledAccount  = NewServiceError("account is not active")
	ErrCreatePwdHash      = NewServiceError("error while hashing password")
	ErrCreateToken        = NewServiceError("error creating token")
	ErrLogout             = NewServiceError("error logout")
	ErrAlreadyExist       = NewServiceError("already exist")
	ErrNothingWasChanged  = NewServiceError("nothing was changed")
	ErrNilMetrics         = NewServiceError("metrics can't be nil")
	ErrNilRepo            = NewServiceError("repository can't be nil")
	ErrEmptyConfig        = NewServiceError("empty config")
	ErrEmptyResult        = NewServiceError("empty result")
)

// FullServiceError возвращает полностью заполненную структуру с типом JointType.
func FullServiceError(message, origin string, initialError error) *errors.BaseError {
	return &errors.BaseError{
		Type:         ErrServiceType,
		Message:      message,
		Origin:       origin,
		InitialError: initialError,
	}
}

// NewServiceError возвращает структуру ошибки с типом JointType и переданным в качестве аргумента сообщением.
func NewServiceError(message string) *errors.BaseError {
	return &errors.BaseError{Type: ErrServiceType, Message: message}
}
