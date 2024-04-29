package postgresql

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const configFilename = "local.yaml"

// dropSchema удаляет использованную схему из БД
func dropSchema(base *PostgreSQL, schema string) {
	if _, err := base.pool.Exec(`DROP SCHEMA IF EXISTS ` + schema + ` CASCADE;`); err != nil {
		log.Print("error dropping schema")
	} else {
		log.Print("schema dropped " + schema)
	}
}

// getConfig возвращает конфигурацию для подключения к БД. Файл с конфигурацией должен лежать в каталоге config.
// Название файла конфигурации указано в константе configFilename в этом тестовом файле
func getConfig() config.PersistentStorage {
	var configPath string
	var err error

	if configPath, err = filepath.Abs("../../../.."); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	if err = os.Setenv("SECURE_CONFIG_PATH", filepath.Join(configPath, "config", configFilename)); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	cfg := config.MustLoad().PersistentStorage
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
		t.Fail()
	}

	dataFromDB, errGet := p.GetAccountLoginData(ctx, data.Login)
	if errGet != nil || data != dataFromDB {
		t.Fail()
	}

	dropSchema(p, p.schema)
	p.Close()
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
