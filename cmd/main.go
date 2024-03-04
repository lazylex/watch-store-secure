package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store/secure/internal/metrics"
	"github.com/lazylex/watch-store/secure/internal/repository/in_memory/redis"
	"github.com/lazylex/watch-store/secure/internal/repository/joint"
	"github.com/lazylex/watch-store/secure/internal/repository/persistent/postgresql"
	"github.com/lazylex/watch-store/secure/internal/service"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustCreate(cfg.Env, cfg.Instance)
	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus, log)
	inMemoryRepo := redis.Create(cfg.Redis)
	persistentRepo := postgresql.Create(cfg.PersistentStorage)
	repo := joint.New(inMemoryRepo, persistentRepo)
	_ = service.New(metrics.Service, repo)

	r := redis.Create(cfg.Redis)

	newId := uuid.New()
	fmt.Println(newId)
	err := r.SaveSession(dto.SessionDTO{Id: newId, Token: "lex", TTL: 200})
	if err != nil {
		fmt.Println(err.Error())
	}

	id, err := r.GetUserUUIDFromSession("lex")
	if err != nil {
		return
	}
	fmt.Println(id.String())

	// TODO удалить, когда будет запущен http-сервер для обработки внешних запросов
	c := make(chan struct{})
	_ = <-c
}
