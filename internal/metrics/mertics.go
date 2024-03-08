package metrics

import (
	"errors"
	"fmt"
	"github.com/lazylex/watch-store/secure/internal/config"
	internalLogger "github.com/lazylex/watch-store/secure/internal/logger"
	"github.com/lazylex/watch-store/secure/internal/ports/metrics/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
)

const NAMESPACE = "secure"

type Metrics struct {
	Service service.MetricsInterface
}

// MustCreate возвращает метрики *Metrics или останавливает программу, если не удалось запустить http сервер для
// работы с Prometheus или занести метрики в регистр
func MustCreate(cfg *config.Prometheus) *Metrics {
	var port = "9323"
	var url = "/metrics"

	if len(cfg.PrometheusPort) > 0 {
		port = cfg.PrometheusPort
	}

	if len(cfg.PrometheusMetricsURL) > 0 {
		url = cfg.PrometheusMetricsURL
	}

	startHTTP(url, port)

	metrics, err := registerMetrics()
	if err != nil {
		slog.With(slog.String(internalLogger.OPLabel, "metrics.MustCreate")).Error(err.Error())
		os.Exit(1)
	}

	return metrics
}

// registerMetrics заносит метрики в регистр и возвращает их. При неудаче возвращает ошибку
func registerMetrics() (*Metrics, error) {
	var (
		err                        error
		loginMetric, authErrMetric *prometheus.CounterVec
	)

	if loginMetric, err = createLoginTotalMetric(); err != nil {
		return nil, err
	}

	if authErrMetric, err = createAuthenticationErrorTotalMetric(); err != nil {
		return nil, err
	}

	return &Metrics{Service: &Service{login: loginMetric, authenticationError: authErrMetric}}, nil
}

// startHTTP запускает http сервер для связи с Prometheus на переданном в функцию порту и url. При неудаче выводит
// ошибку в лог и останавливает программу
func startHTTP(url, port string) {
	go func() {
		mux := http.NewServeMux()
		log := slog.With(internalLogger.OPLabel, "metrics.startHTTP")
		mux.Handle(url, promhttp.Handler())
		log.Info(fmt.Sprintf(":%s%s ready for prometheus", port, url))
		err := http.ListenAndServe(":"+port, mux)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("can't start http server for prometheus")
			os.Exit(1)
		}
	}()
}
