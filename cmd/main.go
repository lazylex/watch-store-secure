package main

import (
	"fmt"
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
	if err := clearScreen(); err != nil {
		slog.Error(err.Error())
	}

	metrics := prometheusMetrics.MustCreate(&cfg.Prometheus)
	inMemoryRepo := redis.MustCreate(cfg.Redis)
	persistentRepo := postgresql.MustCreate(cfg.PersistentStorage)
	repo := joint.New(inMemoryRepo, persistentRepo)
	_ = service.New(metrics.Service, repo)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	fmt.Println() // так красивее, если вывод логов производится в стандартный терминал
	slog.Info(fmt.Sprintf("%s signal received. Shutdown started", sig))
	persistentRepo.Close()
}

func clearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
