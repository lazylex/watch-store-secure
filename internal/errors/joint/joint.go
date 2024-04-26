package joint

import "github.com/lazylex/watch-store/secure/internal/errors"

const jointType = "joint repo"

var (
	ErrDuplicateData      = NewJointError("duplicate data")
	ErrDataNotSaved       = NewJointError("data not saved")
	ErrEmptyResult        = NewJointError("empty result")
	ErrDataTypeConversion = NewJointError("data type conversion failed")
)

type Joint struct {
	errors.BaseError
}

// FullJointError возвращает полностью заполненную структуру Joint
func FullJointError(message, origin string, initialError error) *Joint {
	p := &Joint{}
	p.Type = jointType
	p.Message = message
	p.Origin = origin
	p.InitialError = initialError

	return p
}

// NewJointError возвращает структуру ошибки Joint с переданным в качестве аргумента сообщением
func NewJointError(message string) *Joint {
	p := &Joint{}
	p.Type = jointType
	p.Message = message

	return p
}

// WithOrigin добавляет в структуру место появления ошибки
func (p *Joint) WithOrigin(origin string) *Joint {
	p.Origin = origin
	return p
}
