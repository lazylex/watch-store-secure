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
	stmt := `SELECT uuid, pwd_hash, state FROM accounts WHERE login = $1;`
	row := p.db.QueryRow(stmt, login)
	err := row.Scan(&result.UserId, &result.Hash, &result.State)
	if err != nil {
		return dto.AccountLoginDataDTO{}, err
	}

	return result, nil
}

// SetAccountLoginData сохраняет в БД идентификатор пользователя (сервиса), логин, хеш пароля и состояние учетной записи
func (p *PostgreSQL) SetAccountLoginData(ctx context.Context, data dto.AccountLoginDataDTO) error {
	stmt := `INSERT INTO accounts (uuid, login, pwd_hash, state) values ($1, $2, $3, $4);`
	if _, err := p.db.Exec(stmt, data.UserId, data.Login, data.Hash, data.State); err != nil {
		return err
	}

	return nil
}

// SetAccountState устанавливает состояние учетной записи
func (p *PostgreSQL) SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error {
	stmt := `UPDATE accounts SET state = $1 WHERE login = $2;`
	exec, err := p.db.Exec(stmt, stateDTO.State, stateDTO.Login)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return ErrZeroRowsAffected
	}
	return nil
}

// AddPermission добавляет разрешение в таблицу permissions. В DTO number передавать не обязательно, он вычисляется
// инкременом максимального значения для переданного сервиса. Нумерация для номеров разрешений начинается с нуля
func (p *PostgreSQL) AddPermission(ctx context.Context, perm dto.PermissionDTO) error {
	stmt := `
	INSERT INTO permissions (name, description, service_fk, number)
			VALUES ($1,
			        $2,
			        (SELECT service_id FROM services WHERE name = $3),
			        (SELECT CASE
									WHEN  max(number) IS NULL THEN
									    1
									ELSE 
										max(number) + 1 END
					 			FROM permissions
					 			WHERE service_fk = (SELECT service_id FROM services WHERE name = $3))
					);`
	exec, err := p.db.Exec(stmt, perm.Name, perm.Description, perm.Service)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return ErrZeroRowsAffected
	}
	return nil
}

// AddRole добавляет роль в БД
func (p *PostgreSQL) AddRole(ctx context.Context, data dto.NameAndServiceWithDescriptionDTO) error {
	stmt := `INSERT INTO roles (name, description, service_fk) values ($1, $2, (SELECT service_id FROM services WHERE name=$3));`
	if _, err := p.db.Exec(stmt, data.Name, data.Description, data.Service); err != nil {
		return err
	}

	return nil
}

// AddGroup добавляет группу в БД
func (p *PostgreSQL) AddGroup(ctx context.Context, data dto.NameAndServiceWithDescriptionDTO) error {
	stmt := `INSERT INTO groups (name, description, service_fk) values ($1, $2, (SELECT service_id FROM services WHERE name=$3));`
	if _, err := p.db.Exec(stmt, data.Name, data.Description, data.Service); err != nil {
		return err
	}

	return nil
}

// AddService добавляет сервис в БД
func (p *PostgreSQL) AddService(ctx context.Context, data dto.NameWithDescriptionDTO) error {
	stmt := `INSERT INTO services (name, description) values ($1, $2);`
	if _, err := p.db.Exec(stmt, data.Name, data.Description); err != nil {
		return err
	}

	return nil
}

// AssignPermissionToRole назначает роли разрешение
func (p *PostgreSQL) AssignPermissionToRole(ctx context.Context, data dto.PermissionRoleServiceNamesDTO) error {
	stmt := `INSERT INTO 
    			role_permissions (role_fk, permission_fk)
				VALUES (
                        (SELECT role_id FROM roles WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),
				        (SELECT permission_id FROM permissions WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$3)
				        );`
	if _, err := p.db.Exec(stmt, data.Service, data.Role, data.Permission); err != nil {
		return err
	}
	return nil
}

// AssignRoleToGroup присоединяет роль к группе
func (p *PostgreSQL) AssignRoleToGroup(ctx context.Context, data dto.GroupRoleServiceNamesDTO) error {
	stmt := `INSERT INTO 
    			group_roles (role_fk, groups_fk)
				VALUES (
                        (SELECT role_id FROM roles WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),
				        (SELECT group_id FROM groups WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$3)
				        );`
	if _, err := p.db.Exec(stmt, data.Service, data.Role, data.Group); err != nil {
		return err
	}
	return nil
}
