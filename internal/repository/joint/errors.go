package joint

import (
	"github.com/lazylex/watch-store/secure/internal/errors"
	"github.com/lazylex/watch-store/secure/internal/errors/in_memory"
	"github.com/lazylex/watch-store/secure/internal/errors/joint"
	"github.com/lazylex/watch-store/secure/internal/errors/persistent"
	"strings"
)

const originPlace = "joint → "

// adaptErr переводит пришедшую ошибку к структурированной ошибке Joint.
func adaptErr(err error) error {
	return adaptErrSkipFrames(err, 2)
}

// adaptErrSkipFrames переводит пришедшую ошибку к структурированной ошибке Joint с учетом последовательности
// вызова функций.
func adaptErrSkipFrames(err interface{}, skip int) error {
	if err == nil {
		return nil
	}

	var origin, message string

	if be, ok := err.(*errors.BaseError); ok {
		if len(be.Origin) > 0 {
			origin = be.Origin
		}
		message = be.Message
	}

	if len(origin) == 0 {
		origin = errors.GetFrame(skip).Function
		origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]
	}

	if len(message) == 0 {
		return joint.FullJointError(message, origin, err.(error))
	}

	switch {
	case message == persistent.ErrNoRowsInResultSet.Message:
		return joint.ErrEmptyResult.WithOrigin(origin)
	case message == persistent.ErrZeroRowsAffected.Message:
		return joint.ErrDataNotSaved.WithOrigin(origin)
	case message == persistent.ErrDuplicateKeyValue.Message:
		return joint.ErrDuplicateData.WithOrigin(origin)
	case message == in_memory.ErrNotNumericValue.Message:
		return joint.ErrDataTypeConversion.WithOrigin(origin)
	}

	return joint.FullJointError(message, origin, nil)
}
