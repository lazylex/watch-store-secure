package request_metrics

import (
	"github.com/lazylex/watch-store/secure/internal/adapters/http/router"
	"github.com/lazylex/watch-store/secure/internal/helpers/constants/various"
	"github.com/lazylex/watch-store/secure/internal/metrics"
	"net/http"
	"strings"
	"time"
)

type MiddlewareRequests struct {
	metrics *metrics.Metrics
}

// New конструктор прослойки для http-запросов.
func New(metrics *metrics.Metrics) *MiddlewareRequests {
	return &MiddlewareRequests{metrics: metrics}
}

// BeforeHandle - middleware для метрик, подсчитываемых до выполнения обработчика запроса
func (m *MiddlewareRequests) BeforeHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.requestsInc(r)
		next.ServeHTTP(rw, r)
	})
}

// AfterHandle - middleware для метрик, подсчитываемых после выполнения обработчика запроса.
func (m *MiddlewareRequests) AfterHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := float64(time.Now().UnixMilli()-start.UnixMilli()) * 0.001
			m.metrics.HTTP.RequestsDurationObserve(duration)
		}()

		next.ServeHTTP(rw, r)
	})
}

// requestsInc увеличивает счетчик http-запросов к приложению. Добавляется метка path, содержащая путь запроса.
// Метка path для GET-запросов сохраняется без идентификатора. Если путь не существует, добавляется метка со значением
// "non-existent path".
func (m *MiddlewareRequests) requestsInc(r *http.Request) {
	var path string
	if r.Method == http.MethodGet {
		path = string([]rune(r.RequestURI)[0 : strings.LastIndex(r.RequestURI, "/")+1])
	} else {
		path = r.URL.Path
	}

	if router.IsExistPath(path) {
		m.metrics.HTTP.RequestsTotalInc(map[string]string{metrics.PATH: path})
	} else {
		m.metrics.HTTP.RequestsTotalInc(map[string]string{metrics.PATH: various.NonExistentPath})
	}
}
