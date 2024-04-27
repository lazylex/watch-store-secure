package joint

import "github.com/lazylex/watch-store/secure/internal/errors"

const jointType = "joint repo"

var (
	ErrDuplicateData      = NewJointError("duplicate data")
	ErrDataNotSaved       = NewJointError("data not saved")
	ErrEmptyResult        = NewJointError("empty result")
	ErrDataTypeConversion = NewJointError("data type conversion failed")
	ErrCacheSavedData     = NewJointError("can't save data to cache")
)

// FullJointError возвращает полностью заполненную структуру с типом JointType
func FullJointError(message, origin string, initialError error) *errors.BaseError {
	return &errors.BaseError{
		Type:         jointType,
		Message:      message,
		Origin:       origin,
		InitialError: initialError,
	}
}

// NewJointError возвращает структуру ошибки с типом JointType и переданным в качестве аргумента сообщением
func NewJointError(message string) *errors.BaseError {
	return &errors.BaseError{Type: jointType, Message: message}
}
