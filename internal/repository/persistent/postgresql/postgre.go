package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"github.com/lazylex/watch-store/secure/internal/config"
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"log/slog"
	"os"
)

type PostgreSQL struct {
	db *pgx.Conn
}

var (
	ErrZeroRowsAffected = postgreError("zero rows affected")
)

// postgreError возвращает ошибку с префиксом данного пакета
func postgreError(text string) error {
	return errors.New("postgresql: " + text)
}

// MustCreate возвращает структуру для взаимодействия с базой данных в СУБД PostgreSQL. В случае ошибки завершает
// работу всего приложения
func MustCreate(cfg config.PersistentStorage) *PostgreSQL {
	db, err := pgx.Connect(pgx.ConnConfig{
		Host:     cfg.DatabaseAddress,
		Port:     uint16(cfg.DatabasePort),
		Database: cfg.DatabaseName,
		User:     cfg.DatabaseLogin,
		Password: cfg.DatabasePassword,
	})

	if err != nil {
		exitWithError(err)
	}
	if err = db.Ping(context.Background()); err != nil {
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
	// TODO разобраться, почему не работает плейсхолдер и приходится хардкодить значение в запросе
	stmt := `CREATE TABLE IF NOT EXISTS account (
		uuid varchar(36) NOT NULL UNIQUE,
		login varchar(100) NOT NULL UNIQUE,
		pwd_hash varchar(60) NOT NULL UNIQUE,
		state integer NOT NULL DEFAULT '1',
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

// GetAccountLoginData возвращает необходимые для процесса входа в систему данные пользователя (сервиса)
func (p *PostgreSQL) GetAccountLoginData(ctx context.Context, login loginVO.Login) (dto.AccountLoginDataDTO, error) {
	result := dto.AccountLoginDataDTO{Login: login}
	stmt := `SELECT uuid, pwd_hash, state FROM account WHERE login = $1;`
	row := p.db.QueryRow(stmt, login)
	err := row.Scan(&result.UserId, &result.Hash, &result.State)
	if err != nil {
		return dto.AccountLoginDataDTO{}, err
	}

	return result, nil
}

// SetAccountLoginData сохраняет в БД идентификатор пользователя (сервиса), логин, хеш пароля и состояние учетной записи
func (p *PostgreSQL) SetAccountLoginData(ctx context.Context, data dto.AccountLoginDataDTO) error {
	stmt := `INSERT INTO account (uuid, login, pwd_hash, state) values ($1, $2, $3, $4);`
	if _, err := p.db.Exec(stmt, data.UserId, data.Login, data.Hash, data.State); err != nil {
		return err
	}

	return nil
}

// SetAccountState устанавливает состояние учетной записи
func (p *PostgreSQL) SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error {
	stmt := `UPDATE account SET state = $1 WHERE login = $2;`
	exec, err := p.db.Exec(stmt, stateDTO.State, stateDTO.Login)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return ErrZeroRowsAffected
	}
	return nil
}
