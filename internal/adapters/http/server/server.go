package server

import (
	"context"
	"errors"
	"github.com/lazylex/watch-store/secure/internal/adapters/http/handlers"
	requestMetrics "github.com/lazylex/watch-store/secure/internal/adapters/http/middleware/request_metrics"
	"github.com/lazylex/watch-store/secure/internal/adapters/http/middleware/token_checker"
	"github.com/lazylex/watch-store/secure/internal/adapters/http/router"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/metrics"
	"github.com/lazylex/watch-store/secure/internal/service"
	"log/slog"
	"net/http"
	"os"
)

// Server структура для обработки http-запросов к приложению.
type Server struct {
	cfg     *config.HttpServer // Конфигурация http сервера
	srv     *http.Server       // Структура с параметрами сервера
	mux     *http.ServeMux     // Мультиплексор http запросов
	service *service.Service   // Структура, реализующая логику приложения
}

// MustCreate возвращает готовый к запуску http-сервер (запуск осуществляется функцией MustRun). Если какой-либо из
// переданных параметров равен nil, работа приложения завершается.
func MustCreate(domainService *service.Service, cfg *config.HttpServer, m *metrics.Metrics) *Server {
	if domainService == nil || cfg == nil || m == nil {
		slog.Error("domain service or cfg is nil")
		os.Exit(1)
	}

	mux := http.NewServeMux()
	server := &Server{mux: mux, service: domainService, cfg: cfg}
	server.srv = &http.Server{
		Addr:         server.cfg.Address,
		Handler:      server.mux,
		ReadTimeout:  server.cfg.ReadTimeout,
		WriteTimeout: server.cfg.WriteTimeout,
		IdleTimeout:  server.cfg.IdleTimeout,
	}

	h := handlers.New(domainService, cfg.RequestTimeout)
	router.AssignPathToHandler("/login", server.mux, h.Login)
	router.AssignPathToHandler("/logout", server.mux, h.Logout)
	router.AssignPathToHandler("/get-token", server.mux, h.GetTokenWithPermissions)
	router.AssignPathToHandler("/", server.mux, h.Index)

	tokenMiddleware := token_checker.New(domainService)
	metricsMiddleware := requestMetrics.New(m)

	server.srv.Handler = tokenMiddleware.Checker(server.mux)
	server.srv.Handler = metricsMiddleware.BeforeHandle(server.srv.Handler)
	server.srv.Handler = metricsMiddleware.AfterHandle(server.srv.Handler)
	return server
}

// MustRun производит запуск сервера в отдельной go-рутине. В случае ошибки останавливает работу приложения.
func (s *Server) MustRun() {
	go func() {
		slog.Info("start http server on " + s.srv.Addr)
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server err: startup error. Initial error: " + err.Error())
			os.Exit(1)
		}
	}()
}

// Shutdown производит остановку сервера.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		slog.Error("failed to gracefully shutdown http server")
	} else {
		slog.Info("gracefully shut down http server")
	}
}
