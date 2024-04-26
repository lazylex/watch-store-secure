package postgresql

import (
	"github.com/jackc/pgx"
	"github.com/lazylex/watch-store/secure/internal/errors"
	"github.com/lazylex/watch-store/secure/internal/errors/persistent"
	"strings"
)

const originPlace = "postgresql → "

// processExecResult возвращает ошибку ErrZeroRowsAffected, если при выполнении запроса не было затронуто ни одной
// строки. В противном случае возвращает ошибку, адаптированную под структуру ошибки errors.Persistent
func (p *PostgreSQL) processExecResult(commandTag pgx.CommandTag, err error) error {
	origin := errors.GetFrame(2).Function
	origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]
	if err != nil && strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
		return persistent.ErrDuplicateKeyValue.WithOrigin(origin)
	}

	if commandTag.RowsAffected() == 0 {
		return persistent.ErrZeroRowsAffected
	}

	return adaptErrSkipFrames(err, 3)
}

// adaptErr переводит пришедшую ошибку к структурированной ошибке errors.Persistent
func adaptErr(err error) error {
	return adaptErrSkipFrames(err, 2)
}

// adaptErrSkipFrames переводит пришедшую ошибку к структурированной ошибке errors.Persistent с учетом последовательности
// вызова функций
func adaptErrSkipFrames(err error, skip int) error {
	if err == nil {
		return nil
	}
	origin := errors.GetFrame(skip).Function
	origin = originPlace + origin[strings.LastIndex(origin, ".")+1:]
	if strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
		return persistent.ErrDuplicateKeyValue.WithOrigin(origin)
	}
	if strings.HasPrefix(err.Error(), "no rows in result set") {
		return persistent.ErrNoRowsInResultSet.WithOrigin(origin)
	}

	return persistent.FullPersistentError("", origin, err)
}
