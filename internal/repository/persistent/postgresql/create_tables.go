package postgresql

import "fmt"

// createNotExistedSchemaAndTables создает схему и таблицы в БД, если они отсутствуют.
func (p *PostgreSQL) createNotExistedSchemaAndTables() error {
	var stmt string
	if len(p.schema) > 0 {
		stmt = fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", p.schema)
		if _, err := p.pool.Exec(stmt); err != nil {
			return err
		}
	}

	stmt = `CREATE TABLE IF NOT EXISTS accounts 
		(
			account_id SERIAL PRIMARY KEY,
			uuid UUID NOT NULL UNIQUE,
			login VARCHAR(100) NOT NULL UNIQUE,
			pwd_hash VARCHAR(60) NOT NULL,
			state INTEGER NOT NULL DEFAULT '1'
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS services 
		(
			service_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE, 
			description TEXT
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS instances 
		(
			instance_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			secret VARCHAR(100) NOT NULL,
			service_fk INTEGER NOT NULL REFERENCES services ON DELETE CASCADE
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS permissions
		(
			permission_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			number INTEGER NOT NULL,
			description TEXT,
			service_fk INTEGER NOT NULL REFERENCES services ON DELETE CASCADE,
			UNIQUE (name, service_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS accounts_instances_permissions
		(
			account_fk INTEGER NOT NULL REFERENCES accounts ON DELETE CASCADE,
			instance_fk INTEGER NOT NULL REFERENCES instances ON DELETE CASCADE,
			permission_fk INTEGER NOT NULL REFERENCES permissions ON DELETE CASCADE,
			PRIMARY KEY(account_fk, instance_fk, permission_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS roles 
		(
			role_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			service_fk INTEGER NOT NULL REFERENCES services ON DELETE CASCADE,
			UNIQUE (name, service_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS role_permissions
		(
			role_fk INTEGER NOT NULL REFERENCES roles ON DELETE CASCADE,
			permission_fk INTEGER NOT NULL REFERENCES permissions ON DELETE CASCADE,
			PRIMARY KEY(role_fk, permission_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS account_roles
		(
			role_fk INTEGER NOT NULL REFERENCES roles ON DELETE CASCADE,
			account_fk INTEGER NOT NULL REFERENCES accounts ON DELETE CASCADE,
			PRIMARY KEY(role_fk, account_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS groups 
		(
			group_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			service_fk INTEGER NOT NULL REFERENCES services ON DELETE CASCADE,
			UNIQUE (name, service_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS group_roles
		(
			role_fk INTEGER NOT NULL REFERENCES roles ON DELETE CASCADE,
			group_fk INTEGER NOT NULL REFERENCES groups ON DELETE CASCADE,
			PRIMARY KEY(role_fk, group_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS group_permissions
		(
			permission_fk INTEGER NOT NULL REFERENCES permissions ON DELETE CASCADE,
			group_fk INTEGER NOT NULL REFERENCES groups ON DELETE CASCADE,
			PRIMARY KEY(permission_fk, group_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	stmt = `CREATE TABLE IF NOT EXISTS account_groups
		(
			account_fk INTEGER NOT NULL REFERENCES accounts ON DELETE CASCADE,
			group_fk INTEGER NOT NULL REFERENCES groups ON DELETE CASCADE,
			PRIMARY KEY(account_fk, group_fk)
		)`
	if err := p.createTable(stmt); err != nil {
		return err
	}

	return nil
}

// createTable выполняет переданный в stmt запрос на создание таблицы.
func (p *PostgreSQL) createTable(stmt string) error {
	if _, err := p.pool.Exec(stmt); err != nil {
		return err
	}

	return nil
}
