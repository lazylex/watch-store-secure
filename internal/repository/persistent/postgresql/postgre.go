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

// postgreError возвращает ошибку с префиксом данного пакета
func postgreError(text string) error {
	return errors.New("postgresql: " + text)
}

// MustCreate возвращает структуру для взаимодействия с базой данных в СУБД PostgreSQL. В случае ошибки завершает
// работу всего приложения
func MustCreate(cfg config.PersistentStorage) *PostgreSQL {
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DatabaseAddress, cfg.DatabasePort, cfg.DatabaseLogin, cfg.DatabasePassword, cfg.DatabaseName)

	db, err := sql.Open("postgres", connection)
	if err != nil {
		exitWithError(err)
	}
	if err = db.Ping(); err != nil {
		exitWithError(err)
	}

	slog.Info("successfully ping postgres DB")
	client := &PostgreSQL{db: db}

	if err = client.createNotExistedTables(); err != nil {
		exitWithError(err)
	}

	return client
}

// exitWithError выводит ошибку в лог и завершает приложение
func exitWithError(err error) {
	slog.Error(postgreError(err.Error()).Error())
	os.Exit(1)
}

// createNotExistedTables создает таблицы в БД, если они отсутствуют
func (p *PostgreSQL) createNotExistedTables() error {
	// account table
	stmt := `CREATE TABLE IF NOT EXISTS account (
		uuid varchar(36) NOT NULL UNIQUE,
		login varchar(100) NOT NULL UNIQUE,
		pwd_hash varchar(60) NOT NULL UNIQUE,
		enabled integer NOT NULL DEFAULT '1',
		PRIMARY KEY (uuid))`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	}

	// TODO add tables
	slog.Debug(postgreError("only one table created in DB").Error())

	return nil
}

// Close закрывает соединение с БД
func (p *PostgreSQL) Close() {
	if err := p.db.Close(); err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("closed postgres db connection")
	}
}

// GetUserIdAndPasswordHash возвращает идентификатор пользователя и хэш пароля для пользователя с переданным логином
func (p *PostgreSQL) GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error) {
	// TODO implement
	slog.Debug(postgreError("GetUserIdAndPasswordHash not implemented").Error())
	return dto.UserIdWithPasswordHashDTO{}, nil
}
