package main

import (
	"fmt"
	"github.com/lazylex/watch-store/secure/internal/adapters/http/server"
	"github.com/lazylex/watch-store/secure/internal/adapters/message_broker/kafka"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/logger"
	prometheusMetrics "github.com/lazylex/watch-store/secure/internal/metrics"
	"github.com/lazylex/watch-store/secure/internal/repository/in_memory/redis"
	"github.com/lazylex/watch-store/secure/internal/repository/joint"
	"github.com/lazylex/watch-store/secure/internal/repository/persistent/postgresql"
	"github.com/lazylex/watch-store/secure/internal/service"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	cfg := config.MustLoad()
	slog.SetDefault(logger.MustCreate(cfg.Env, cfg.Instance))
	clearScreen()

	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus)
	inMemoryRepo := redis.MustCreate(cfg.Redis, cfg.TTL)
	persistentRepo := postgresql.MustCreate(cfg.PersistentStorage)
	repo := joint.MustCreate(inMemoryRepo, persistentRepo)
	domainService := service.MustCreate(metrics.Service, &repo, cfg.Secure)

	httpServer := server.MustCreate(domainService, &cfg.HttpServer)
	httpServer.MustRun()

	if cfg.UseKafka {
		kafka.MustRun(&cfg.Kafka)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	fmt.Println() // так красивее, если вывод логов производится в стандартный терминал
	slog.Info(fmt.Sprintf("%s signal received. Shutdown started", sig))

	httpServer.Shutdown()
	persistentRepo.Close()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
