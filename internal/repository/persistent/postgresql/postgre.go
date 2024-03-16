package postgresql

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"log/slog"
)

type PostgreSQL struct {
}

func postgreError(text string) error {
	return errors.New("postgre: " + text)
}

func Create(cfg config.PersistentStorage) *PostgreSQL {
	// TODO implement
	return &PostgreSQL{}
}

func (p *PostgreSQL) GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error) {
	// TODO implement
	slog.Debug(postgreError("GetUserIdAndPasswordHash not implemented").Error())
	return dto.UserIdWithPasswordHashDTO{}, nil
}
