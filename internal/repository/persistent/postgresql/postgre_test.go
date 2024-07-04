package postgresql

import (
	"context"
	"errors"
	"github.com/google/uuid"
	storageConfig "github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/errors/persistent"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const configFilename = "local.yaml"

var baseConfig storageConfig.PersistentStorage

// testConfig возвращает конфигурацию для подключения к БД. Файл с конфигурацией должен лежать в каталоге testConfig.
// Название файла конфигурации указано в константе configFilename в этом тестовом файле
func testConfig() storageConfig.PersistentStorage {
	var configPath string
	var err error
	var cfg storageConfig.PersistentStorage

	if len(baseConfig.DatabaseAddress) == 0 {
		if configPath, err = filepath.Abs("../../../.."); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		if err = os.Setenv("SECURE_CONFIG_PATH", filepath.Join(configPath, "config", configFilename)); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		cfg = storageConfig.MustLoad().PersistentStorage
		baseConfig = cfg
	} else {
		cfg = baseConfig
	}

	return cfg
}

// postgreSQL возвращает ссылку на готовую для работы с БД структуру PostgreSQL
func postgreSQL(t *testing.T) *PostgreSQL {
	cfg := testConfig()
	p := MustCreateForTest(cfg)
	t.Cleanup(p.DropCurrentTestSchema)
	return p
}

func TestPostgreSQL_SetAndGetAccountLoginData(t *testing.T) {
	var err error
	p := postgreSQL(t)
	ctx := context.Background()

	data := dto.UserIdLoginHashState{
		Login:  "test_user",
		UserId: uuid.New(),
		Hash:   "$2a$14$qXnQ8n9U0FItXkto3Sf8XuvZny48y4iZLTluWZtZszTrc7REdzUAy",
		State:  account_state.Enabled,
	}

	err = p.SetAccountLoginData(ctx, &data)
	if err != nil {
		t.Fatal()
	}

	dataFromDB, errGet := p.AccountLoginData(ctx, data.Login)
	if errGet != nil || data != dataFromDB {
		t.Fail()
	}
}

func TestPostgreSQL_ErrGetAccountLoginData(t *testing.T) {
	p := postgreSQL(t)

	if _, err := p.AccountLoginData(context.Background(), "non-existent user"); !errors.Is(err, persistent.ErrNoRowsInResultSet) {
		t.Fail()
	}
}

func TestPostgreSQL_ErrCreateConnection(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		cfg := testConfig()
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
	p := postgreSQL(t)
	data := dto.NameDescription{
		Name:        "test_service",
		Description: "just test service",
	}
	if p.CreateService(context.Background(), &data) != nil {
		t.Fail()
	}
}

func TestPostgreSQL_ErrCreateDuplicateService(t *testing.T) {
	p := postgreSQL(t)
	data := dto.NameDescription{
		Name:        "test_service",
		Description: "just test service",
	}
	if p.CreateService(context.Background(), &data) != nil {
		t.Fatal()
	}
	if !errors.Is(p.CreateService(context.Background(), &data), persistent.ErrDuplicateKeyValue) {
		t.Fail()
	}
}

func TestPostgreSQL_BigTest(t *testing.T) {
	p := postgreSQL(t)
	ctx := context.Background()

	if p.CreateService(ctx, &dto.NameDescription{Name: "service1", Description: "description 1"}) != nil {
		t.Fatal()
	}

	if p.CreateGroup(ctx, &dto.NameServiceDescription{Name: "group1", Description: "group1 description", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreateRole(ctx, &dto.NameServiceDescription{Name: "role1", Description: "r1", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreatePermission(ctx, &dto.NameServiceDescription{Name: "perm1", Description: "p1", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreatePermission(ctx, &dto.NameServiceDescription{Name: "perm2", Description: "p2", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.CreatePermission(ctx, &dto.NameServiceDescription{Name: "perm3", Description: "p3", Service: "service1"}) != nil {
		t.Fatal()
	}

	userId := uuid.New()
	if p.SetAccountLoginData(ctx, &dto.UserIdLoginHashState{
		Login:  "test_user",
		UserId: userId,
		Hash:   "$2a$14$qXnQ8n9U0FItXkto3Sf8XuvZny48y4iZLTluWZtZszTrc7REdzUAy",
		State:  account_state.Enabled,
	}) != nil {
		t.Fatal()
	}

	if p.AssignPermissionToGroup(ctx, &dto.GroupPermissionService{
		Group:      "group1",
		Permission: "perm1",
		Service:    "service1",
	}) != nil {
		t.Fatal()
	}

	if p.AssignPermissionToRole(ctx, &dto.PermissionRoleService{
		Permission: "perm2",
		Role:       "role1",
		Service:    "service1",
	}) != nil {
		t.Fatal()
	}

	if p.AssignRoleToGroup(ctx, &dto.GroupRoleService{
		Group:   "group1",
		Role:    "role1",
		Service: "service1",
	}) != nil {
		t.Fatal()
	}

	if p.AssignGroupToAccount(ctx, &dto.UserIdGroupService{
		UserId:  userId,
		Group:   "group1",
		Service: "service1",
	}) != nil {
		t.Fatal()
	}

	if permissions, err := p.ServicePermissionsForAccount(ctx, &dto.UserIdService{
		UserId:  userId,
		Service: "service1",
	}); err != nil {
		t.Fatal()
	} else {
		if len(permissions) != 2 {
			t.Fatal("not enough permissions count")
		}
	}

	if perm, err := p.InstancePermissionsForAccount(ctx, &dto.UserIdInstance{
		UserId:   userId,
		Instance: "instance1",
	}); len(perm) > 0 || err != nil {
		t.Fatal()
	}

	if p.CreateOrUpdateInstance(ctx, &dto.NameServiceSecret{
		Name:    "instance1",
		Service: "service1",
		Secret:  "нет никакого секрета",
	}) != nil {
		t.Fatal()
	}

	if p.AssignInstancePermissionToAccount(ctx, &dto.UserIdInstancePermission{
		UserId:     userId,
		Instance:   "instance1",
		Permission: "perm3",
	}) != nil {
		t.Fatal()
	}

	if perm, err := p.InstancePermissionsForAccount(ctx, &dto.UserIdInstance{
		UserId:   userId,
		Instance: "instance1",
	}); len(perm) != 1 || err != nil {
		t.Fatal()
	}

	if numbers, err := p.InstancePermissionsNumbersForAccount(ctx, &dto.UserIdInstance{
		UserId:   userId,
		Instance: "instance1",
	}); err != nil || len(numbers) != 1 || numbers[0] != 3 {
		t.Fatal()
	}

	if numbers, err := p.ServicePermissionsNumbersForAccount(ctx, &dto.UserIdService{
		UserId:  userId,
		Service: "service1",
	}); err != nil || len(numbers) != 2 {
		t.Fatal()
	}

	if p.SetAccountState(ctx, &dto.LoginState{Login: "test_user", State: account_state.Disabled}) != nil {
		t.Fatal()
	}

	if logins, err := p.AccountsLoginsByState(ctx, account_state.Enabled); err != nil || len(logins) != 0 {
		t.Fatal()
	}

	if logins, err := p.AccountsLoginsByState(ctx, account_state.Disabled); err != nil || len(logins) != 1 || logins[0] != "test_user" {
		t.Fatal()
	}

	if p.MaxConnections() != p.maxConnections {
		slog.Error("Вообще бесполезная проверка, но так покрытие тестами полнее")
		t.Fatal()
	}

	if number, err := p.PermissionNumber(ctx, "perm3", "instance1"); err != nil || number != 3 {
		t.Fatal()
	}

	if p.CreateRole(ctx, &dto.NameServiceDescription{Name: "role2", Description: "r2", Service: "service1"}) != nil {
		t.Fatal()
	}

	if p.AssignRoleToAccount(ctx, &dto.UserIdRoleService{
		UserId:  userId,
		Role:    "role2",
		Service: "service1",
	}) != nil {
		t.Fatal()
	}
}
