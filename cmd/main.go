package main

import (
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store/secure/internal/metrics"
	"github.com/lazylex/watch-store/secure/internal/repository/in_memory/redis"
	"github.com/lazylex/watch-store/secure/internal/repository/joint"
	"github.com/lazylex/watch-store/secure/internal/repository/persistent/postgresql"
	"github.com/lazylex/watch-store/secure/internal/service"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	slog.SetDefault(logger.MustCreate(cfg.Env, cfg.Instance))
	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus)
	inMemoryRepo := redis.MustCreate(cfg.Redis)
	persistentRepo := postgresql.Create(cfg.PersistentStorage)
	repo := joint.New(inMemoryRepo, persistentRepo)
	_ = service.New(metrics.Service, repo)

	// TODO удалить, когда будет запущен http-сервер для обработки внешних запросов
	c := make(chan struct{})
	_ = <-c
}
