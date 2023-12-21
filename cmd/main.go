package main

import (
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store/secure/internal/metrics"
	"github.com/lazylex/watch-store/secure/internal/service"
)

func main() {
	cfg := config.MustLoad()
	log := logger.MustCreate(cfg.Env, cfg.Instance)
	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus, log)
	_ = service.New(metrics.Service)

	// TODO удалить, когда будет запущен http-сервер для обработки внешних запросов
	c := make(chan struct{})
	_ = <-c
}
