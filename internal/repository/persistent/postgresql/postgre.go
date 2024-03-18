package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

type PostgreSQL struct {
	db *sql.DB
}

func postgreError(text string) error {
	return errors.New("postgre: " + text)
}

func MustCreate(cfg config.PersistentStorage) *PostgreSQL {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseAddress, cfg.DatabasePort, cfg.DatabaseLogin, cfg.DatabasePassword, cfg.DatabaseName)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if err = db.Ping(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	slog.Info("successfully ping postgres DB")
	return &PostgreSQL{db: db}
}

func (p *PostgreSQL) Close() {
	if err := p.db.Close(); err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("closed postgres db connection")
	}
}

func (p *PostgreSQL) GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error) {
	// TODO implement
	slog.Debug(postgreError("GetUserIdAndPasswordHash not implemented").Error())
	return dto.UserIdWithPasswordHashDTO{}, nil
}
