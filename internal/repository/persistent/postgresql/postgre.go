package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	loginVO "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"log/slog"
	"math/rand"
	"os"
	"strings"
)

const testSchemaPrefix = "test_schema_"

type PostgreSQL struct {
	pool           *pgx.ConnPool
	maxConnections int
	schema         string
}

// MustCreate возвращает структуру для взаимодействия с базой данных в СУБД PostgreSQL. В случае ошибки завершает
// работу всего приложения
func MustCreate(cfg config.PersistentStorage) *PostgreSQL {
	schema := "public"
	if len(cfg.DatabaseSchema) > 0 {
		schema = cfg.DatabaseSchema
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:          cfg.DatabaseAddress,
			Port:          uint16(cfg.DatabasePort),
			Database:      cfg.DatabaseName,
			User:          cfg.DatabaseLogin,
			Password:      cfg.DatabasePassword,
			RuntimeParams: map[string]string{"search_path": schema},
		},
		MaxConnections: cfg.DatabaseMaxOpenConnections,
	})

	if err != nil {
		slog.Error(adaptErr(err).Error())
		os.Exit(1)
	} else {
		slog.Info("successfully create connection poll to postgres DB")
	}

	client := &PostgreSQL{pool: pool, maxConnections: cfg.DatabaseMaxOpenConnections, schema: schema}

	if err = client.createNotExistedSchemaAndTables(); err != nil {
		slog.Error(adaptErr(err).Error())
		os.Exit(1)
	}

	return client
}

// MustCreateForTest возвращает структуру для взаимодействия с тестовой базой данных в СУБД PostgreSQL. Переданная в
// конфигурации схема игнорируется и меняется на сгенерированную случайным образом (префикс из константы
// testSchemaPrefix и случайное целое число). В остальном идентично функции MustCreate
func MustCreateForTest(cfg config.PersistentStorage) *PostgreSQL {
	cfg.DatabaseSchema = fmt.Sprintf("%s%d", testSchemaPrefix, rand.Int())
	return MustCreate(cfg)
}

// Close закрывает пул соединений с БД
func (p *PostgreSQL) Close() {
	p.pool.Close()
	slog.Info("closed postgres pool")
}

// DropCurrentTestSchema удаляет текущую схему со всеми данными, если она предназначалась для тестов. Это определяется
// по суффиксу схемы, который должен быть test_schema_ для тестовых схем
func (p *PostgreSQL) DropCurrentTestSchema() {
	if !strings.HasPrefix(p.schema, testSchemaPrefix) {
		slog.Warn(fmt.Sprintf("can't dropping current schema. It's not start with %s prefix", testSchemaPrefix))
		return
	}

	if _, err := p.pool.Exec(fmt.Sprintf("DROP SCHEMA %s CASCADE", p.schema)); err == nil {
		p.Close()
	} else {
		slog.Error(adaptErr(err).Error())
	}
}

func (p *PostgreSQL) GetMaxConnections() int {
	return p.maxConnections
}

// GetAccountLoginData возвращает необходимые для процесса входа в систему данные пользователя (сервиса)
func (p *PostgreSQL) GetAccountLoginData(ctx context.Context, login loginVO.Login) (dto.AccountLoginDataDTO, error) {
	result := dto.AccountLoginDataDTO{Login: login}
	stmt := `SELECT uuid, pwd_hash, state FROM accounts WHERE login = $1;`
	row := p.pool.QueryRowEx(ctx, stmt, nil, login)
	err := row.Scan(&result.UserId, &result.Hash, &result.State)
	if err != nil {
		return dto.AccountLoginDataDTO{}, adaptErr(err)
	}

	return result, nil
}

// SetAccountLoginData сохраняет в БД идентификатор пользователя (сервиса), логин, хеш пароля и состояние учетной записи
func (p *PostgreSQL) SetAccountLoginData(ctx context.Context, data *dto.AccountLoginDataDTO) error {
	stmt := `INSERT INTO accounts (uuid, login, pwd_hash, state) values ($1, $2, $3, $4);`
	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.UserId, data.Login, data.Hash, data.State))
}

// SetAccountState устанавливает состояние учетной записи
func (p *PostgreSQL) SetAccountState(ctx context.Context, data *dto.LoginStateDTO) error {
	stmt := `UPDATE accounts SET state = $1 WHERE login = $2;`
	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.State, data.Login))
}

// CreatePermission добавляет разрешение в таблицу permissions
func (p *PostgreSQL) CreatePermission(ctx context.Context, data *dto.PermissionWithoutNumberDTO) error {
	stmt := `INSERT INTO permissions (name, description, service_fk, number)
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

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Name, data.Description, data.Service))
}

// CreateRole добавляет роль в БД
func (p *PostgreSQL) CreateRole(ctx context.Context, data *dto.NameAndServiceWithDescriptionDTO) error {
	stmt := `INSERT INTO roles (name, description, service_fk)
			VALUES ($1, $2, (SELECT service_id FROM services WHERE name=$3));`
	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Name, data.Description, data.Service))
}

// CreateGroup добавляет группу в БД
func (p *PostgreSQL) CreateGroup(ctx context.Context, data *dto.NameAndServiceWithDescriptionDTO) error {
	stmt := `INSERT INTO groups (name, description, service_fk)
			VALUES ($1, $2, (SELECT service_id FROM services WHERE name=$3));`
	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Name, data.Description, data.Service))
}

// CreateService добавляет сервис в БД
func (p *PostgreSQL) CreateService(ctx context.Context, data *dto.NameWithDescriptionDTO) error {
	stmt := `INSERT INTO services (name, description) VALUES ($1, $2);`
	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Name, data.Description))
}

// CreateInstance добавляет в БД название экземпляра сервиса
func (p *PostgreSQL) CreateInstance(ctx context.Context, data *dto.NameAndServiceDTO) error {
	stmt := `INSERT INTO instances (name, service_fk)
			VALUES ($1,
			        (SELECT service_id FROM services WHERE name=$2));`
	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Name, data.Service))
}

// AssignPermissionToRole назначает роли разрешение
func (p *PostgreSQL) AssignPermissionToRole(ctx context.Context, data *dto.PermissionRoleServiceNamesDTO) error {
	stmt := `INSERT INTO role_permissions (role_fk, permission_fk)
			VALUES (
            	(SELECT role_id 
				FROM roles 
				WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),

				(SELECT permission_id 
				FROM permissions 
				WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$3)
			);`

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Service, data.Role, data.Permission))
}

// AssignRoleToGroup присоединяет роль к группе
func (p *PostgreSQL) AssignRoleToGroup(ctx context.Context, data *dto.GroupRoleServiceNamesDTO) error {
	stmt := `INSERT INTO group_roles (role_fk, group_fk)
			VALUES (
				(SELECT role_id 
				FROM roles 
				WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),

				(SELECT group_id 
				FROM groups 
				WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$3)
			);`

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Service, data.Role, data.Group))
}

// AssignRoleToAccount назначает роль учетной записи
func (p *PostgreSQL) AssignRoleToAccount(ctx context.Context, data *dto.RoleServiceNamesWithUserIdDTO) error {
	stmt := `INSERT INTO account_roles (role_fk, account_fk)
			VALUES (
				(SELECT role_id 
                FROM roles 
                WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),

				(SELECT account_id
				FROM accounts
				WHERE uuid = $3) 
			);`

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Service, data.Role, data.UserId))
}

// AssignGroupToAccount назначает группу учетной записи
func (p *PostgreSQL) AssignGroupToAccount(ctx context.Context, data *dto.GroupServiceNamesWithUserIdDTO) error {
	stmt := `INSERT INTO account_groups (group_fk, account_fk)
			VALUES (
				(SELECT group_id 
                FROM groups 
                WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),

				(SELECT account_id
				FROM accounts
				WHERE uuid = $3) 
			);`

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Service, data.Group, data.UserId))
}

// AssignInstancePermissionToAccount прикрепляет разрешение конкретного экземпляра сервиса к учетной записи
func (p *PostgreSQL) AssignInstancePermissionToAccount(ctx context.Context, data *dto.InstanceAndPermissionNamesWithUserIdDTO) error {
	cte := `WITH instance_cte AS (SELECT instance_id, service_fk
                      FROM instances
                      WHERE name = $1)`
	stmt := cte + `
		INSERT
		INTO accounts_instances_permissions (account_fk, instance_fk, permission_fk)
		VALUES ((SELECT account_id
				 FROM accounts
				 WHERE uuid = $2),
				(SELECT instance_id FROM instance_cte),
				(SELECT permission_id
				 FROM permissions
				 WHERE name = $3
				   AND service_fk IN (SELECT service_fk FROM instance_cte)));`

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Instance, data.UserId, data.Permission))
}

// AssignPermissionToGroup назначает разрешения группе
func (p *PostgreSQL) AssignPermissionToGroup(ctx context.Context, data *dto.GroupPermissionServiceNamesDTO) error {
	stmt := `INSERT INTO group_permissions (group_fk, permission_fk)
			VALUES (
	       	(SELECT group_id
				FROM groups
				WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$2),
	
				(SELECT permission_id
				FROM permissions
				WHERE service_fk = (SELECT service_id FROM services WHERE name=$1) AND name=$3)
			);`

	return p.processExecResult(p.pool.ExecEx(ctx, stmt, nil, data.Service, data.Group, data.Permission))
}

// GetInstancePermissionsForAccount возвращает название, номер и описание разрешений аккаунта для экземпляра сервиса
func (p *PostgreSQL) GetInstancePermissionsForAccount(ctx context.Context, data *dto.InstanceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error) {
	cte := `
	WITH account_cte AS (SELECT account_id
                     FROM accounts
                     WHERE uuid = $1)`

	stmt := cte + `
	SELECT name,
		   number,
		   description
	FROM permissions
	WHERE permission_id IN
		  (SELECT permission_fk
		   FROM accounts_instances_permissions
		   WHERE account_fk = (SELECT account_id FROM account_cte)
	
			 AND instance_fk IN
				 (SELECT instance_id
				  FROM instances
				  WHERE name = $2))

	ORDER BY number`

	rows, err := p.pool.QueryEx(ctx, stmt, nil, data.UserId, data.Instance)
	defer rows.Close()

	if err != nil {
		return nil, adaptErr(err)
	}

	result := make([]dto.PermissionWithoutServiceDTO, 0)
	var number int
	var name, description string

	for rows.Next() {
		if err = rows.Scan(&name, &number, &description); err != nil {
			return result, adaptErr(err)
		}
		result = append(result, dto.PermissionWithoutServiceDTO{Name: name, Number: number, Description: description})

	}

	return result, nil
}

// GetInstancePermissionsNumbersForAccount возвращает номера разрешений аккаунта для экземпляра сервиса
func (p *PostgreSQL) GetInstancePermissionsNumbersForAccount(ctx context.Context, data *dto.InstanceNameWithUserIdDTO) ([]int, error) {
	cte := `
	WITH account_cte AS (SELECT account_id
                     FROM accounts
                     WHERE uuid = $1)`

	stmt := cte + `
	SELECT number
	FROM permissions
	WHERE permission_id IN
		  (SELECT permission_fk
		   FROM accounts_instances_permissions
		   WHERE account_fk = (SELECT account_id FROM account_cte)
	
			 AND instance_fk IN
				 (SELECT instance_id
				  FROM instances
				  WHERE name = $2))

	ORDER BY number`

	rows, err := p.pool.QueryEx(ctx, stmt, nil, data.UserId, data.Instance)
	defer rows.Close()
	if err != nil {
		return nil, adaptErr(err)
	}

	result := make([]int, 0)
	var number int

	for rows.Next() {
		if err = rows.Scan(&number); err != nil {
			return result, adaptErr(err)
		}
		result = append(result, number)

	}

	return result, nil
}

// GetServicePermissionsForAccount возвращает название, номер и описание разрешений аккаунта для сервиса (без разрешений
// для экземпляра)
func (p *PostgreSQL) GetServicePermissionsForAccount(ctx context.Context, data *dto.ServiceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error) {
	cte := `
	WITH account_cte AS (SELECT account_id
                     FROM accounts
                     WHERE uuid = $1),
     groups_cte AS (SELECT group_fk
                    FROM account_groups
                    WHERE account_fk = (SELECT account_id FROM account_cte)),
     service_cte AS (SELECT service_id
                     FROM services
                     WHERE name = $2)`

	stmt := cte + `
	SELECT name,
		   number,
		   description
	FROM permissions
	WHERE permission_id IN
		  (SELECT permission_fk
		   FROM role_permissions
		   WHERE role_fk IN
				 (SELECT role_fk
				  FROM group_roles
				  WHERE group_fk IN (SELECT group_fk FROM groups_cte)
	
				  UNION
	
				  SELECT role_fk
				  FROM account_roles
				  WHERE account_fk = (SELECT account_id FROM account_cte))
	
		   UNION
	
		   SELECT permission_fk
		   FROM group_permissions
		   WHERE group_fk IN (SELECT group_fk FROM groups_cte)
		   )
	
	  AND service_fk = (SELECT service_id FROM service_cte)
	ORDER BY number`

	rows, err := p.pool.QueryEx(ctx, stmt, nil, data.UserId, data.Service)
	defer rows.Close()

	if err != nil {
		return nil, adaptErr(err)
	}

	result := make([]dto.PermissionWithoutServiceDTO, 0)
	var number int
	var name, description string

	for rows.Next() {
		if err = rows.Scan(&name, &number, &description); err != nil {
			return result, adaptErr(err)
		}
		result = append(result, dto.PermissionWithoutServiceDTO{Name: name, Number: number, Description: description})

	}

	return result, nil
}

// GetServicePermissionsNumbersForAccount возвращает номера разрешений аккаунта для сервиса (без разрешений для
// экземпляра)
func (p *PostgreSQL) GetServicePermissionsNumbersForAccount(ctx context.Context, data *dto.ServiceNameWithUserIdDTO) ([]int, error) {
	cte := `
	WITH account_cte AS (SELECT account_id
                     FROM accounts
                     WHERE uuid = $1),
     groups_cte AS (SELECT group_fk
                    FROM account_groups
                    WHERE account_fk = (SELECT account_id FROM account_cte)),
     service_cte AS (SELECT service_id
                     FROM services
                     WHERE name = $2)`

	stmt := cte + `
	SELECT number
	FROM permissions
	WHERE permission_id IN
		  (SELECT permission_fk
		   FROM role_permissions
		   WHERE role_fk IN
				 (SELECT role_fk
				  FROM group_roles
				  WHERE group_fk IN (SELECT group_fk FROM groups_cte)
	
				  UNION
	
				  SELECT role_fk
				  FROM account_roles
				  WHERE account_fk = (SELECT account_id FROM account_cte))
	
		   UNION
	
		   SELECT permission_fk
		   FROM group_permissions
		   WHERE group_fk IN (SELECT group_fk FROM groups_cte)
		   )
	
	  AND service_fk = (SELECT service_id FROM service_cte)
	ORDER BY number`

	rows, err := p.pool.QueryEx(ctx, stmt, nil, data.UserId, data.Service)
	defer rows.Close()

	if err != nil {
		return nil, adaptErr(err)
	}

	result := make([]int, 0)
	var number int

	for rows.Next() {
		if err = rows.Scan(&number); err != nil {
			return result, adaptErr(err)
		}
		result = append(result, number)

	}

	return result, nil
}

// GetPermissionNumber возвращает номер разрешения для заданного экземпляра сервиса
func (p *PostgreSQL) GetPermissionNumber(ctx context.Context, name, instance string) (int, error) {
	var number int
	stmt := `SELECT number FROM permissions WHERE name=$1 AND service_fk=(SELECT service_fk FROM instances WHERE name=$2)`
	row := p.pool.QueryRowEx(ctx, stmt, nil, name, instance)
	if err := row.Scan(&number); err != nil {
		return 0, adaptErr(err)
	}

	return number, nil
}

// GetAccountsLoginsByState возвращает список логинов пользователей с переданным функции состоянием
func (p *PostgreSQL) GetAccountsLoginsByState(ctx context.Context, state account_state.State) ([]loginVO.Login, error) {
	stmt := `SELECT login FROM accounts WHERE state = $1`
	rows, err := p.pool.QueryEx(ctx, stmt, nil, state)
	defer rows.Close()

	if err != nil {
		return nil, adaptErr(err)
	}

	result := make([]loginVO.Login, 0)

	var accLogin loginVO.Login

	for rows.Next() {
		if err = rows.Scan(&accLogin); err != nil {
			return result, adaptErr(err)
		}
		result = append(result, accLogin)

	}

	return result, nil
}
