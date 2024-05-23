package metrics

import (
	"github.com/lazylex/watch-store/secure/internal/adapters/http/router"
	"github.com/lazylex/watch-store/secure/internal/helpers/constants/various"
	"github.com/prometheus/client_golang/prometheus"
)

type HTTP struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

// RequestsTotalInc инкремент счетчика запросов.
func (h *HTTP) RequestsTotalInc(labels map[string]string) {
	h.requests.With(labels).Inc()
}

// RequestsDurationObserve внесение данных о длительности запроса.
func (h *HTTP) RequestsDurationObserve(duration float64) {
	h.duration.With(prometheus.Labels{}).Observe(duration)
}

// createHTTPRequestDurationSecondsBucketMetric создает и регистрирует метрику http_request_duration_seconds_bucket.
func createHTTPRequestDurationSecondsBucketMetric() (*prometheus.HistogramVec, error) {
	var err error
	requestDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: NAMESPACE,
		Name:      "http_request_duration_seconds_bucket",
		Help:      "duration of the request",
	}, []string{})
	if err = prometheus.Register(requestDuration); err != nil {
		return nil, err
	}

	requestDuration.With(prometheus.Labels{})

	return requestDuration, nil
}

// createHTTPRequestsTotalMetric создает и регистрирует метрику http_requests_total, являющуюся счетчиком http-запросов.
func createHTTPRequestsTotalMetric() (*prometheus.CounterVec, error) {
	var err error
	requests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "http_requests_total",
		Namespace: NAMESPACE,
		Help:      "Count of http requests",
	}, []string{PATH})
	if err = prometheus.Register(requests); err != nil {
		return nil, err
	}

	for _, path := range router.ExistentPaths() {
		requests.With(prometheus.Labels{PATH: path})
	}
	requests.With(prometheus.Labels{PATH: various.NonExistentPath})

	return requests, nil
}
