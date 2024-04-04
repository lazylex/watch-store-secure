package postgresql

import (
	"log/slog"
)

// createNotExistedTables создает таблицы в БД, если они отсутствуют
func (p *PostgreSQL) createNotExistedTables() error {
	// TODO добавить каскадное удаление
	var stmt string
	// accounts table
	stmt = `CREATE TABLE IF NOT EXISTS accounts 
		(
			account_id SERIAL PRIMARY KEY,
			uuid UUID NOT NULL UNIQUE,
			login VARCHAR(100) NOT NULL UNIQUE,
			pwd_hash VARCHAR(60) NOT NULL,
			state INTEGER NOT NULL DEFAULT '1'
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "accounts"))
	}

	// services table
	stmt = `CREATE TABLE IF NOT EXISTS services 
		(
			service_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE, 
			description TEXT
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "services"))
	}

	// instances table
	stmt = `CREATE TABLE IF NOT EXISTS instances 
		(
			instance_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			service_fk INTEGER NOT NULL REFERENCES services
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "instances"))
	}

	// permissions table
	stmt = `CREATE TABLE IF NOT EXISTS permissions
		(
			permission_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			number INTEGER NOT NULL,
			description TEXT,
			service_fk INTEGER NOT NULL REFERENCES services
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "permissions"))
	}

	// accounts instances permissions table
	stmt = `CREATE TABLE IF NOT EXISTS accounts_instances_permissions
		(
			account_fk INTEGER NOT NULL REFERENCES accounts,
			instance_fk INTEGER NOT NULL REFERENCES instances,
			permission_fk INTEGER NOT NULL REFERENCES permissions,
			PRIMARY KEY(account_fk, instance_fk, permission_fk)
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "accounts_instances_permissions"))
	}

	// roles table
	stmt = `CREATE TABLE IF NOT EXISTS roles 
		(
			role_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "roles"))
	}

	// role permissions table
	stmt = `CREATE TABLE IF NOT EXISTS role_permissions
		(
			role_fk INTEGER NOT NULL REFERENCES roles,
			permission_fk INTEGER NOT NULL REFERENCES permissions,
			PRIMARY KEY(role_fk, permission_fk)
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "role_permissions"))
	}

	// account roles table
	stmt = `CREATE TABLE IF NOT EXISTS account_roles
		(
			role_fk INTEGER NOT NULL REFERENCES roles,
			account_fk INTEGER NOT NULL REFERENCES accounts,
			PRIMARY KEY(role_fk, account_fk)
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "account_roles"))
	}

	// groups table
	stmt = `CREATE TABLE IF NOT EXISTS groups 
		(
			group_id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "groups"))
	}

	// group roles table
	stmt = `CREATE TABLE IF NOT EXISTS group_roles
		(
			role_fk INTEGER NOT NULL REFERENCES roles,
			groups_fk INTEGER NOT NULL REFERENCES groups,
			PRIMARY KEY(role_fk, groups_fk)
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "group_roles"))
	}

	// group permissions table
	stmt = `CREATE TABLE IF NOT EXISTS group_permissions
		(
			permission_fk INTEGER NOT NULL REFERENCES permissions,
			groups_fk INTEGER NOT NULL REFERENCES groups,
			PRIMARY KEY(permission_fk, groups_fk)
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "group_permissions"))
	}

	// account groups table
	stmt = `CREATE TABLE IF NOT EXISTS account_groups
		(
			account_fk INTEGER NOT NULL REFERENCES accounts,
			groups_fk INTEGER NOT NULL REFERENCES groups,
			PRIMARY KEY(account_fk, groups_fk)
		)`
	if _, err := p.db.Exec(stmt); err != nil {
		return err
	} else {
		slog.Info("created table", slog.String("table name", "account_groups"))
	}

	return nil
}
