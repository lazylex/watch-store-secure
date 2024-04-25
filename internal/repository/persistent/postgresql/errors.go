package postgresql

import (
	"github.com/jackc/pgx"
	"github.com/lazylex/watch-store/secure/internal/lexerr"
	"strings"
)

// processExecResult возвращает ошибку ErrZeroRowsAffected, если при выполнении запроса не было затронуто ни одной
// строки. В противном случае возвращает ошибку, адаптированную под структуру ошибки lexerr.Persistent
func (p *PostgreSQL) processExecResult(commandTag pgx.CommandTag, err error) error {
	origin := lexerr.GetFrame(2).Function
	origin = origin[strings.LastIndex(origin, ".")+1:]
	if err != nil && strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
		return lexerr.ErrDuplicateKeyValue.WithOrigin(origin)
	}

	if commandTag.RowsAffected() == 0 {
		return lexerr.ErrZeroRowsAffected
	}

	return adaptErrSkipFrames(err, 3)
}

// adaptErr переводит пришедшую ошибку к структурированной ошибке lexerr.Persistent
func adaptErr(err error) error {
	return adaptErrSkipFrames(err, 2)
}

// adaptErrSkipFrames переводит пришедшую ошибку к структурированной ошибке lexerr.Persistent с учетом последовательности
// вызова функций
func adaptErrSkipFrames(err error, skip int) error {
	if err == nil {
		return nil
	}
	origin := lexerr.GetFrame(skip).Function
	origin = origin[strings.LastIndex(origin, ".")+1:]
	if strings.HasPrefix(err.Error(), "ERROR: duplicate key value violates unique constraint") {
		return lexerr.ErrDuplicateKeyValue.WithOrigin(origin)
	}
	if strings.HasPrefix(err.Error(), "no rows in result set") {
		return lexerr.ErrNoRowsInResultSet.WithOrigin(origin)
	}

	return lexerr.FullPersistentError("", origin, err)
}
