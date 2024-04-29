package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/errors/persistent"
	"log/slog"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const configFilename = "local.yaml"

var baseConfig config.PersistentStorage

// dropSchema удаляет использованную схему из БД
func dropSchema(base *PostgreSQL, schema string) {
	_, _ = base.pool.Exec(`DROP SCHEMA IF EXISTS ` + schema + ` CASCADE;`)
}

// getConfig возвращает конфигурацию для подключения к БД. Файл с конфигурацией должен лежать в каталоге config.
// Название файла конфигурации указано в константе configFilename в этом тестовом файле
func getConfig() config.PersistentStorage {
	var configPath string
	var err error
	var cfg config.PersistentStorage

	if len(baseConfig.DatabaseAddress) == 0 {
		if configPath, err = filepath.Abs("../../../.."); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if err = os.Setenv("SECURE_CONFIG_PATH", filepath.Join(configPath, "config", configFilename)); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		cfg = config.MustLoad().PersistentStorage
		baseConfig = cfg
	} else {
		cfg = baseConfig
	}

	cfg.DatabaseSchema = fmt.Sprintf("test_schema_%d", rand.Int())

	return cfg
}

// getPostgreSQL возвращает ссылку на готовую для работы с БД структуру PostgreSQL
func getPostgreSQL() *PostgreSQL {
	cfg := getConfig()
	return MustCreate(cfg)
}

func TestPostgreSQL_SetAndGetAccountLoginData(t *testing.T) {
	var err error
	p := getPostgreSQL()
	ctx := context.Background()

	data := dto.AccountLoginDataDTO{
		Login:  "test_user",
		UserId: uuid.New(),
		Hash:   "$2a$14$qXnQ8n9U0FItXkto3Sf8XuvZny48y4iZLTluWZtZszTrc7REdzUAy",
		State:  account_state.Enabled,
	}

	err = p.SetAccountLoginData(ctx, data)
	if err != nil {
		t.Fatal()
	}

	dataFromDB, errGet := p.GetAccountLoginData(ctx, data.Login)
	if errGet != nil || data != dataFromDB {
		t.Fail()
	}

	dropSchema(p, p.schema)
	p.Close()
}

func TestPostgreSQL_ErrGetAccountLoginData(t *testing.T) {
	p := getPostgreSQL()
	defer dropSchema(p, p.schema)
	defer p.Close()
	if _, err := p.GetAccountLoginData(context.Background(), "non-existent user"); !errors.Is(err, persistent.ErrNoRowsInResultSet) {
		t.Fail()
	}
}

func TestPostgreSQL_ErrCreateConnection(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		cfg := getConfig()
		cfg.DatabaseName = ""
		MustCreate(cfg)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestPostgreSQL_ErrCreateConnection")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestPostgreSQL_CreateService(t *testing.T) {
	p := getPostgreSQL()
	defer dropSchema(p, p.schema)
	defer p.Close()
	data := dto.NameWithDescriptionDTO{
		Name:        "test_service",
		Description: "just test service",
	}
	if p.CreateService(context.Background(), data) != nil {
		t.Fail()
	}
}

func TestPostgreSQL_ErrCreateDuplicateService(t *testing.T) {
	p := getPostgreSQL()
	defer dropSchema(p, p.schema)
	defer p.Close()
	data := dto.NameWithDescriptionDTO{
		Name:        "test_service",
		Description: "just test service",
	}
	if p.CreateService(context.Background(), data) != nil {
		t.Fatal()
	}
	if !errors.Is(p.CreateService(context.Background(), data), persistent.ErrDuplicateKeyValue) {
		t.Fail()
	}
}

func TestPostgreSQL_BigTest(t *testing.T) {
	p := getPostgreSQL()
	defer dropSchema(p, p.schema)
	defer p.Close()
	ctx := context.Background()

	if p.CreateService(ctx, dto.NameWithDescriptionDTO{Name: "service1", Description: "description 1"}) != nil {
		t.Fatal()
	}

	if p.CreateGroup(ctx, dto.NameAndServiceWithDescriptionDTO{Name: "group1", Description: "group1 description", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreateRole(ctx, dto.NameAndServiceWithDescriptionDTO{Name: "role1", Description: "r1", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreatePermission(ctx, dto.PermissionWithoutNumberDTO{Name: "perm1", Description: "p1", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreatePermission(ctx, dto.PermissionWithoutNumberDTO{Name: "perm2", Description: "p2", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreatePermission(ctx, dto.PermissionWithoutNumberDTO{Name: "perm3", Description: "p3", Service: "service1"}) != nil {
		t.Fatal()
	}

	userId := uuid.New()
	if p.SetAccountLoginData(ctx, dto.AccountLoginDataDTO{
		Login:  "test_user",
		UserId: userId,
		Hash:   "$2a$14$qXnQ8n9U0FItXkto3Sf8XuvZny48y4iZLTluWZtZszTrc7REdzUAy",
		State:  account_state.Enabled,
	}) != nil {
		t.Fatal()
	}

	if p.AssignPermissionToGroup(ctx, dto.GroupPermissionServiceNamesDTO{
		Group:      "group1",
		Permission: "perm1",
		Service:    "service1",
	}) != nil {
		t.Fatal()
	}

	if p.AssignPermissionToRole(ctx, dto.PermissionRoleServiceNamesDTO{
		Permission: "perm2",
		Role:       "role1",
		Service:    "service1",
	}) != nil {
		t.Fatal()
	}

	if p.AssignRoleToGroup(ctx, dto.GroupRoleServiceNamesDTO{
		Group:   "group1",
		Role:    "role1",
		Service: "service1",
	}) != nil {
		t.Fatal()
	}

	// TODO добавить perm3 к экземпляру сервиса service1 и попробовать его считать

	if p.AssignGroupToAccount(ctx, dto.GroupServiceNamesWithUserIdDTO{
		UserId:  userId,
		Group:   "group1",
		Service: "service1",
	}) != nil {
		t.Fatal()
	}

	if permissions, err := p.GetServicePermissionsForAccount(ctx, dto.ServiceNameWithUserIdDTO{
		UserId:  userId,
		Service: "service1",
	}); err != nil {
		t.Fatal()
	} else {
		if len(permissions) != 2 {
			t.Fatal("not enough permissions count")
		}
	}

	if perm, err := p.GetInstancePermissionsForAccount(ctx, dto.InstanceNameWithUserIdDTO{
		UserId:   userId,
		Instance: "instance1",
	}); len(perm) > 0 || err != nil {
		t.Fatal()
	}
}
